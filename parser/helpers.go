package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	headerRe    = regexp.MustCompile(`(?i)^(flowchart|graph)\s+(\w+)`)
	subgraphRe  = regexp.MustCompile(`(?i)^subgraph\s+(.*)$`)
	pipeLabelRe = regexp.MustCompile(
		`^(?P<left>.+?)\s*(?P<arrow><[-.=ox]*[-=]+[-.=ox]*>|<[-.=ox]*[-=]+|[-.=ox]*[-=]+>|[-.=ox]*[-=]+)\|(?P<label>.+?)\|\s*(?P<right>.+)$`,
	)
	quotedLabelArrowRe = regexp.MustCompile(
		`^(?P<left>.+?)\s*(?P<start><)?(?P<dash1>[-.=ox]*[-=]+[-.=ox]*)\s+"(?P<label>[^"]+)"\s+(?P<dash2>[-.=ox]*[-=]+[-.=ox]*)(?P<end>>)?\s*(?P<right>.+)$`,
	)
	labelArrowRe = regexp.MustCompile(
		`^(?P<left>.+?)\s*(?P<start><)?(?P<dash1>[-.=ox]*[-=]+[-.=ox]*)\s+(?P<label>[^<>=]+?)\s+(?P<dash2>[-.=ox]*[-=]+[-.=ox]*)(?P<end>>)?\s*(?P<right>.+)$`,
	)
	compactDottedLabelRe = regexp.MustCompile(
		`^(?P<left>.+?)\s*(?P<start><)?(?P<dash1>[-=ox]*[-=]+[-=ox]*)\.(?P<label>[^<>=|].*?)\.(?P<dash2>[-.=ox]*[-=]+[-.=ox]*)(?P<end>>)?\s*(?P<right>.+)$`,
	)
	arrowRe = regexp.MustCompile(
		`^(?P<left>.+?)\s*(?P<arrow><[-.=ox]*[-=]+[-.=ox]*>|<[-.=ox]*[-=]+|[-.=ox]*[-=]+>|[-.=ox]*[-=]+)\s*(?P<right>.+)$`,
	)
	arrowTokenRe = regexp.MustCompile(
		`<[-.=ox]*[-=]+[-.=ox]*>|<[-.=ox]*[-=]+|[-.=ox]*[-=]+>|[-.=ox]*[-=]+`,
	)
)

// edgeMeta holds parsed metadata about an edge arrow.
type edgeMeta struct {
	directed        bool
	arrowStart      bool
	arrowEnd        bool
	arrowStartKind  *ir.EdgeArrowhead
	arrowEndKind    *ir.EdgeArrowhead
	startDecoration *ir.EdgeDecoration
	endDecoration   *ir.EdgeDecoration
	style           ir.EdgeStyle
}

// stripTrailingComment removes %% comments from the end of a line,
// respecting quoted strings.
func stripTrailingComment(line string) string {
	var quote rune
	var out strings.Builder
	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		ch := runes[i]
		if quote != 0 {
			if ch == quote {
				quote = 0
			}
			out.WriteRune(ch)
			continue
		}
		if ch == '"' || ch == '\'' {
			quote = ch
			out.WriteRune(ch)
			continue
		}
		if ch == '%' && i+1 < len(runes) && runes[i+1] == '%' {
			break
		}
		out.WriteRune(ch)
	}
	return strings.TrimSpace(out.String())
}

// splitStatements splits a line on ; at depth 0 (outside brackets/quotes).
func splitStatements(line string) []string {
	var parts []string
	var current strings.Builder
	depth := 0
	var quote rune
	escaped := false

	for _, ch := range line {
		if escaped {
			current.WriteRune(ch)
			escaped = false
			continue
		}
		if ch == '\\' {
			current.WriteRune(ch)
			escaped = true
			continue
		}
		if quote != 0 {
			if ch == quote {
				quote = 0
			}
			current.WriteRune(ch)
			continue
		}
		if ch == '"' || ch == '\'' {
			quote = ch
			current.WriteRune(ch)
			continue
		}
		switch ch {
		case '[', '(', '{':
			depth++
			current.WriteRune(ch)
		case ']', ')', '}':
			if depth > 0 {
				depth--
			}
			current.WriteRune(ch)
		case ';':
			if depth == 0 {
				trimmed := strings.TrimSpace(current.String())
				if trimmed != "" {
					parts = append(parts, trimmed)
				}
				current.Reset()
			} else {
				current.WriteRune(ch)
			}
		default:
			current.WriteRune(ch)
		}
	}

	trimmed := strings.TrimSpace(current.String())
	if trimmed != "" {
		parts = append(parts, trimmed)
	}
	return parts
}

// maskBracketContent replaces characters inside [], (), {}, "", ” with spaces
// while preserving byte positions. This prevents edge-detection regexes from
// matching dashes inside labels.
func maskBracketContent(line string) string {
	runes := []rune(line)
	result := make([]rune, len(runes))
	depthSquare := 0
	depthParen := 0
	depthCurly := 0
	inDoubleQuote := false
	inSingleQuote := false
	var prevChar rune

	for i, ch := range runes {
		inBracket := depthSquare > 0 || depthParen > 0 || depthCurly > 0
		inQuote := inDoubleQuote || inSingleQuote

		switch {
		case ch == '[' && !inQuote:
			depthSquare++
			result[i] = ch
		case ch == ']' && !inQuote && depthSquare > 0:
			depthSquare--
			result[i] = ch
		case ch == '(' && !inQuote && !inBracket:
			depthParen++
			result[i] = ch
		case ch == ')' && !inQuote && depthParen > 0:
			depthParen--
			result[i] = ch
		case ch == '{' && !inQuote && !inBracket:
			depthCurly++
			result[i] = ch
		case ch == '}' && !inQuote && depthCurly > 0:
			depthCurly--
			result[i] = ch
		case ch == '"' && prevChar != '\\':
			inDoubleQuote = !inDoubleQuote
			if inBracket || inQuote {
				result[i] = ' '
			} else {
				result[i] = ch
			}
		case ch == '\'' && prevChar != '\\':
			inSingleQuote = !inSingleQuote
			if inBracket || inQuote {
				result[i] = ' '
			} else {
				result[i] = ch
			}
		default:
			if inBracket || inQuote {
				result[i] = ' '
			} else {
				result[i] = ch
			}
		}
		prevChar = ch
	}
	return string(result)
}

// splitEdgeChain detects A-->B-->C chains (2+ arrows with no labels) and
// splits them into individual edge statements.
func splitEdgeChain(line string) []string {
	masked := maskBracketContent(line)

	// If any label pattern matches, this is not a simple chain.
	if pipeLabelRe.MatchString(masked) ||
		quotedLabelArrowRe.MatchString(line) ||
		labelArrowRe.MatchString(masked) ||
		compactDottedLabelRe.MatchString(masked) {
		return nil
	}

	matches := arrowTokenRe.FindAllStringIndex(masked, -1)
	if len(matches) < 2 {
		return nil
	}

	nodes := make([]string, 0, len(matches)+1)
	arrows := make([]string, 0, len(matches))
	lastIdx := 0

	for _, m := range matches {
		nodes = append(nodes, strings.TrimSpace(line[lastIdx:m[0]]))
		arrows = append(arrows, strings.TrimSpace(line[m[0]:m[1]]))
		lastIdx = m[1]
	}
	nodes = append(nodes, strings.TrimSpace(line[lastIdx:]))

	if len(nodes) != len(arrows)+1 {
		return nil
	}

	// Attach leading pipe labels to the preceding arrow.
	for i := 1; i < len(nodes); i++ {
		trimmed := strings.TrimLeft(nodes[i], " \t")
		if strings.HasPrefix(trimmed, "|") {
			endIdx := strings.Index(trimmed[1:], "|")
			if endIdx >= 0 {
				labelLen := endIdx + 2
				label := trimmed[:labelLen]
				rest := strings.TrimLeft(trimmed[labelLen:], " \t")
				arrows[i-1] += label
				nodes[i] = rest
			}
		}
	}

	for _, node := range nodes {
		if node == "" {
			return nil
		}
	}

	statements := make([]string, 0, len(arrows))
	for i := 0; i < len(arrows); i++ {
		statements = append(statements, nodes[i]+" "+arrows[i]+" "+nodes[i+1])
	}
	return statements
}

// parseEdgeLine tries each regex pattern in priority order to parse an edge
// line into (left, label, right, edgeMeta).
func parseEdgeLine(line string) (left string, label *string, right string, meta edgeMeta, ok bool) {
	masked := maskBracketContent(line)

	// Helper to extract from original line using match positions from masked line.
	namedSub := func(re *regexp.Regexp, s string) map[string]string {
		match := re.FindStringSubmatch(s)
		if match == nil {
			return nil
		}
		result := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		return result
	}

	namedIdx := func(re *regexp.Regexp, s string) map[string][2]int {
		idx := re.FindStringSubmatchIndex(s)
		if idx == nil {
			return nil
		}
		result := make(map[string][2]int)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" && idx[2*i] >= 0 {
				result[name] = [2]int{idx[2*i], idx[2*i+1]}
			}
		}
		return result
	}

	// 1. Pipe label: A -->|label| B
	if idx := namedIdx(pipeLabelRe, masked); idx != nil {
		leftRange := idx["left"]
		rightRange := idx["right"]
		labelRange := idx["label"]
		arrowRange := idx["arrow"]

		l := strings.TrimSpace(line[leftRange[0]:leftRange[1]])
		r := strings.TrimSpace(line[rightRange[0]:rightRange[1]])
		lbl := strings.TrimSpace(line[labelRange[0]:labelRange[1]])
		arrow := strings.TrimSpace(line[arrowRange[0]:arrowRange[1]])

		if lbl != "" && l != "" && r != "" {
			m := parseEdgeMeta(arrow)
			return l, &lbl, r, m, true
		}
	}

	// 2. Quoted label: A -- "label" --> B (match on original line)
	if caps := namedSub(quotedLabelArrowRe, line); caps != nil {
		l := strings.TrimSpace(caps["left"])
		r := strings.TrimSpace(caps["right"])
		lbl := strings.TrimSpace(caps["label"])
		if lbl != "" && l != "" && r != "" {
			start := caps["start"]
			dash1 := caps["dash1"]
			dash2 := caps["dash2"]
			end := caps["end"]
			arrow := start + dash1 + dash2 + end
			m := parseEdgeMeta(arrow)
			return l, &lbl, r, m, true
		}
	}

	// 3. Label arrow: A -- label --> B
	if idx := namedIdx(labelArrowRe, masked); idx != nil {
		leftRange := idx["left"]
		rightRange := idx["right"]
		labelRange := idx["label"]

		l := strings.TrimSpace(line[leftRange[0]:leftRange[1]])
		r := strings.TrimSpace(line[rightRange[0]:rightRange[1]])
		lblRaw := strings.TrimSpace(line[labelRange[0]:labelRange[1]])
		lbl := strings.Trim(lblRaw, "|")
		lbl = strings.TrimSpace(lbl)

		if lbl != "" && l != "" && r != "" {
			start := ""
			dash1 := ""
			dash2 := ""
			end := ""
			if rng, has := idx["start"]; has {
				start = masked[rng[0]:rng[1]]
			}
			if rng, has := idx["dash1"]; has {
				dash1 = masked[rng[0]:rng[1]]
			}
			if rng, has := idx["dash2"]; has {
				dash2 = masked[rng[0]:rng[1]]
			}
			if rng, has := idx["end"]; has {
				end = masked[rng[0]:rng[1]]
			}
			arrow := start + dash1 + dash2 + end
			m := parseEdgeMeta(arrow)
			return l, &lbl, r, m, true
		}
	}

	// 4. Compact dotted label: A -.label.-> B
	if idx := namedIdx(compactDottedLabelRe, masked); idx != nil {
		leftRange := idx["left"]
		rightRange := idx["right"]
		labelRange := idx["label"]

		l := strings.TrimSpace(line[leftRange[0]:leftRange[1]])
		r := strings.TrimSpace(line[rightRange[0]:rightRange[1]])
		lbl := strings.Trim(strings.TrimSpace(line[labelRange[0]:labelRange[1]]), ".")

		if lbl != "" && l != "" && r != "" {
			start := ""
			dash1 := ""
			dash2 := ""
			end := ""
			if rng, has := idx["start"]; has {
				start = masked[rng[0]:rng[1]]
			}
			if rng, has := idx["dash1"]; has {
				dash1 = masked[rng[0]:rng[1]]
			}
			if rng, has := idx["dash2"]; has {
				dash2 = masked[rng[0]:rng[1]]
			}
			if rng, has := idx["end"]; has {
				end = masked[rng[0]:rng[1]]
			}
			arrow := start + dash1 + "." + dash2 + end
			m := parseEdgeMeta(arrow)
			return l, &lbl, r, m, true
		}
	}

	// 5. Simple arrow: A --> B
	idx := namedIdx(arrowRe, masked)
	if idx == nil {
		return "", nil, "", edgeMeta{}, false
	}

	leftRange := idx["left"]
	rightRange := idx["right"]
	arrowRange := idx["arrow"]

	l := strings.TrimSpace(line[leftRange[0]:leftRange[1]])
	arrow := strings.TrimSpace(masked[arrowRange[0]:arrowRange[1]])
	r := strings.TrimSpace(line[rightRange[0]:rightRange[1]])

	// Check for leading decoration on right side (e.g., "o B" after arrow).
	if dec, rest, found := extractLeadingDecoration(r); found {
		arrow += string(dec)
		r = rest
	}

	if l == "" || r == "" || arrow == "" {
		return "", nil, "", edgeMeta{}, false
	}

	// Check for trailing pipe label on the right side: |label| node
	var lbl *string
	rightToken := r
	if strings.HasPrefix(r, "|") {
		if endIdx := strings.Index(r[1:], "|"); endIdx >= 0 {
			labelStr := strings.TrimSpace(r[1 : endIdx+1])
			rest := strings.TrimSpace(r[endIdx+2:])
			if rest != "" {
				lbl = &labelStr
				rightToken = rest
			}
		}
	}

	if rightToken == "" {
		return "", nil, "", edgeMeta{}, false
	}

	m := parseEdgeMeta(arrow)
	return l, lbl, rightToken, m, true
}

// extractLeadingDecoration checks if right starts with 'o' or 'x' followed by
// whitespace, which indicates a decoration character was split from the arrow.
func extractLeadingDecoration(right string) (rune, string, bool) {
	runes := []rune(right)
	if len(runes) < 2 {
		return 0, "", false
	}
	first := runes[0]
	if first != 'o' && first != 'x' {
		return 0, "", false
	}
	rest := string(runes[1:])
	if len(rest) > 0 && (rest[0] == ' ' || rest[0] == '\t') {
		return first, strings.TrimLeft(rest, " \t"), true
	}
	return 0, "", false
}

// parseEdgeMeta extracts style, direction, and decorations from an arrow string.
func parseEdgeMeta(arrow string) edgeMeta {
	trimmed := strings.TrimSpace(arrow)
	var startDecoration *ir.EdgeDecoration
	var endDecoration *ir.EdgeDecoration

	runes := []rune(trimmed)
	if len(runes) > 0 {
		switch runes[0] {
		case 'o':
			dec := ir.DecCircle
			startDecoration = &dec
			runes = runes[1:]
		case 'x':
			dec := ir.DecCross
			startDecoration = &dec
			runes = runes[1:]
		}
	}

	if len(runes) > 0 {
		switch last := runes[len(runes)-1]; last {
		case 'o':
			dec := ir.DecCircle
			endDecoration = &dec
			runes = runes[:len(runes)-1]
		case 'x':
			dec := ir.DecCross
			endDecoration = &dec
			runes = runes[:len(runes)-1]
		}
	}

	trimmed = string(runes)
	arrowStart := strings.HasPrefix(trimmed, "<")
	arrowEnd := strings.HasSuffix(trimmed, ">")

	var style ir.EdgeStyle
	if strings.Contains(trimmed, "=") {
		style = ir.Thick
	} else if strings.Contains(trimmed, ".") {
		style = ir.Dotted
	} else {
		style = ir.Solid
	}

	directed := arrowStart || arrowEnd

	return edgeMeta{
		directed:        directed,
		arrowStart:      arrowStart,
		arrowEnd:        arrowEnd,
		startDecoration: startDecoration,
		endDecoration:   endDecoration,
		style:           style,
	}
}

// parseNodeToken parses a node token like "A[Start]" into its components.
func parseNodeToken(token string) (id string, label *string, shape *ir.NodeShape, classes []string) {
	base, classes := splitInlineClasses(token)
	trimmed := strings.TrimSpace(base)

	// Try asymmetric shape first: A>label]
	if asymID, asymLabel, asymShape, found := splitAsymmetricLabel(trimmed); found {
		return asymID, &asymLabel, &asymShape, classes
	}

	// Try bracket-based shapes.
	if splitID, splitLabel, splitShape, found := splitIDLabel(trimmed); found {
		return splitID, &splitLabel, &splitShape, classes
	}

	// Bare node ID (first whitespace-delimited token).
	parts := strings.Fields(trimmed)
	if len(parts) > 0 {
		id = parts[0]
	}
	return id, nil, nil, classes
}

// splitInlineClasses splits "A[label]:::className" into ("A[label]", ["className"]).
func splitInlineClasses(token string) (string, []string) {
	parts := strings.Split(token, ":::")
	base := strings.TrimSpace(parts[0])
	var classes []string
	for _, part := range parts[1:] {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			classes = append(classes, trimmed)
		}
	}
	return base, classes
}

// splitAsymmetricLabel parses the A>label] shape.
func splitAsymmetricLabel(token string) (string, string, ir.NodeShape, bool) {
	trimmed := strings.TrimSpace(token)
	if strings.Contains(trimmed, "[") {
		return "", "", 0, false
	}
	pos := strings.Index(trimmed, ">")
	if pos < 0 {
		return "", "", 0, false
	}
	if !strings.HasSuffix(trimmed, "]") {
		return "", "", 0, false
	}
	id := strings.TrimSpace(trimmed[:pos])
	if id == "" {
		return "", "", 0, false
	}
	label := strings.TrimSpace(trimmed[pos+1 : len(trimmed)-1])
	if label == "" {
		return "", "", 0, false
	}
	label = stripQuotes(label)
	return id, label, ir.Asymmetric, true
}

// splitIDLabel detects the shape from bracket patterns in the token.
func splitIDLabel(token string) (string, string, ir.NodeShape, bool) {
	// Try [...] shapes
	if start := strings.Index(token, "["); start >= 0 && strings.HasSuffix(token, "]") {
		id := strings.TrimSpace(token[:start])
		if id != "" {
			raw := token[start:]
			label, shape := parseShapeFromBrackets(raw)
			return id, label, shape, true
		}
	}

	// Try (...) shapes
	if start := strings.Index(token, "("); start >= 0 && strings.HasSuffix(token, ")") {
		id := strings.TrimSpace(token[:start])
		if id != "" {
			raw := token[start:]
			label, shape := parseShapeFromParens(raw)
			return id, label, shape, true
		}
	}

	// Try {...} shapes
	if start := strings.Index(token, "{"); start >= 0 && strings.HasSuffix(token, "}") {
		id := strings.TrimSpace(token[:start])
		if id != "" {
			raw := token[start:]
			label, shape := parseShapeFromBraces(raw)
			return id, label, shape, true
		}
	}

	return "", "", 0, false
}

// parseShapeFromBrackets parses [...] variants.
func parseShapeFromBrackets(raw string) (string, ir.NodeShape) {
	trimmed := strings.TrimSpace(raw)

	if strings.HasPrefix(trimmed, "[/") && strings.HasSuffix(trimmed, "/]") {
		return stripQuotes(trimmed[2 : len(trimmed)-2]), ir.Parallelogram
	}
	if strings.HasPrefix(trimmed, "[\\") && strings.HasSuffix(trimmed, "\\]") {
		return stripQuotes(trimmed[2 : len(trimmed)-2]), ir.ParallelogramAlt
	}
	if strings.HasPrefix(trimmed, "[/") && strings.HasSuffix(trimmed, "\\]") {
		return stripQuotes(trimmed[2 : len(trimmed)-2]), ir.Trapezoid
	}
	if strings.HasPrefix(trimmed, "[\\") && strings.HasSuffix(trimmed, "/]") {
		return stripQuotes(trimmed[2 : len(trimmed)-2]), ir.TrapezoidAlt
	}
	if strings.HasPrefix(trimmed, "[[") && strings.HasSuffix(trimmed, "]]") {
		return stripQuotes(trimmed[2 : len(trimmed)-2]), ir.Subroutine
	}
	if strings.HasPrefix(trimmed, "[(") && strings.HasSuffix(trimmed, ")]") {
		return stripQuotes(trimmed[2 : len(trimmed)-2]), ir.Cylinder
	}
	if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
		inner := trimmed[1 : len(trimmed)-1]
		if strings.HasPrefix(inner, "(") && strings.HasSuffix(inner, ")") {
			return stripQuotes(inner[1 : len(inner)-1]), ir.Stadium
		}
		return stripQuotes(inner), ir.Rectangle
	}
	return stripQuotes(trimmed), ir.Rectangle
}

// parseShapeFromParens parses (...) variants.
func parseShapeFromParens(raw string) (string, ir.NodeShape) {
	trimmed := strings.TrimSpace(raw)

	if strings.HasPrefix(trimmed, "(((") && strings.HasSuffix(trimmed, ")))") {
		return stripQuotes(trimmed[3 : len(trimmed)-3]), ir.DoubleCircle
	}
	if strings.HasPrefix(trimmed, "((") && strings.HasSuffix(trimmed, "))") {
		return stripQuotes(trimmed[2 : len(trimmed)-2]), ir.DoubleCircle
	}
	if strings.HasPrefix(trimmed, "(") && strings.HasSuffix(trimmed, ")") {
		inner := trimmed[1 : len(trimmed)-1]
		if strings.HasPrefix(inner, "[") && strings.HasSuffix(inner, "]") {
			return stripQuotes(inner[1 : len(inner)-1]), ir.Stadium
		}
		return stripQuotes(inner), ir.RoundRect
	}
	return stripQuotes(trimmed), ir.RoundRect
}

// parseShapeFromBraces parses {...} variants.
func parseShapeFromBraces(raw string) (string, ir.NodeShape) {
	trimmed := strings.TrimSpace(raw)

	if strings.HasPrefix(trimmed, "{{") && strings.HasSuffix(trimmed, "}}") {
		return stripQuotes(trimmed[2 : len(trimmed)-2]), ir.Hexagon
	}
	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		return stripQuotes(trimmed[1 : len(trimmed)-1]), ir.Diamond
	}
	return stripQuotes(trimmed), ir.Diamond
}

// stripQuotes removes surrounding double or single quotes from a string.
func stripQuotes(input string) string {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) >= 2 {
		if (trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"') ||
			(trimmed[0] == '\'' && trimmed[len(trimmed)-1] == '\'') {
			return trimmed[1 : len(trimmed)-1]
		}
	}
	return trimmed
}

// parseSubgraphHeader parses the "subgraph rest..." part after the keyword.
func parseSubgraphHeader(input string) (id *string, label string, classes []string) {
	base, classes := splitInlineClasses(input)
	trimmed := strings.TrimSpace(base)

	if trimmed == "" {
		label = "Subgraph"
		return nil, label, classes
	}

	// Try split_id_label for patterns like sg1[Group]
	if splitID, splitLabel, _, found := splitIDLabel(trimmed); found {
		return &splitID, splitLabel, classes
	}

	// If no quotes and single token, use as both ID and label.
	if !strings.ContainsAny(trimmed, "\"'") {
		parts := strings.Fields(trimmed)
		if len(parts) == 1 {
			token := parts[0]
			return &token, token, classes
		}
	}

	// Use as label, no ID.
	label = stripQuotes(trimmed)
	return nil, label, classes
}

// parseDirectionLine checks if a line is "direction LR" etc.
func parseDirectionLine(line string) (ir.Direction, bool) {
	parts := strings.Fields(line)
	if len(parts) == 2 && parts[0] == "direction" {
		return ir.DirectionFromToken(parts[1])
	}
	return ir.TopDown, false
}

// parseNodeOnly attempts to parse a line as a standalone node declaration.
// Returns false if the line contains "--" (likely an edge).
func parseNodeOnly(line string) (id string, label *string, shape *ir.NodeShape, classes []string, ok bool) {
	if strings.Contains(line, "--") {
		return "", nil, nil, nil, false
	}
	id, label, shape, classes = parseNodeToken(line)
	if id == "" {
		return "", nil, nil, nil, false
	}
	return id, label, shape, classes, true
}

// extractQuotedText extracts the text between the first pair of double quotes.
// Returns empty string if no quoted text is found.
func extractQuotedText(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' {
		if end := strings.Index(s[1:], "\""); end >= 0 {
			return s[1 : end+1]
		}
	}
	return ""
}

// splitAndTrimCommas splits a string by commas, trims whitespace and quotes.
func splitAndTrimCommas(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, "\"")
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
