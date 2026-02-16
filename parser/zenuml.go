package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

// zenBlockKind identifies what opened a curly-brace block.
type zenBlockKind int

const (
	zenMsgBlock     zenBlockKind = iota // A.method() { ... }
	zenIfBlock                          // if(cond) {
	zenElseIfBlock                      // else if(cond) {
	zenElseBlock                        // else {
	zenLoopBlock                        // while/for/forEach/loop {
	zenTryBlock                         // try {
	zenCatchBlock                       // catch {
	zenFinallyBlock                     // finally {
	zenOptBlock                         // opt {
	zenParBlock                         // par {
	zenGroupBlock                       // group Name {
)

// zenBlock represents an open curly-brace block on the parsing stack.
type zenBlock struct {
	kind   zenBlockKind
	caller string // caller context before this block opened
	target string // for zenMsgBlock: the activated participant
}

var (
	zenStarterRe  = regexp.MustCompile(`(?i)^@Starter\s*\(\s*(\w+)\s*\)$`)
	zenAnnotRe    = regexp.MustCompile(`(?i)^@(Actor|Boundary|Control|Entity|Database|Collections|Queue)\s+(\w+)(?:\s+as\s+(.+?))?$`)
	zenAliasRe    = regexp.MustCompile(`^(\w+)\s+as\s+(.+)$`)
	zenGroupRe    = regexp.MustCompile(`(?i)^group\s+(.+?)\s*\{$`)
	zenIfRe       = regexp.MustCompile(`^if\s*\((.+)\)\s*\{$`)
	zenElseIfRe   = regexp.MustCompile(`^else\s+if\s*\((.+)\)\s*\{$`)
	zenLoopRe     = regexp.MustCompile(`(?i)^(while|for|forEach|loop)\s*(?:\((.+?)\))?\s*\{$`)
	zenTryRe      = regexp.MustCompile(`(?i)^try\s*(?:\(\))?\s*\{$`)
	zenCatchRe    = regexp.MustCompile(`(?i)^catch\s*(?:\(\))?\s*\{$`)
	zenFinallyRe  = regexp.MustCompile(`(?i)^finally\s*(?:\(\))?\s*\{$`)
	zenOptRe      = regexp.MustCompile(`(?i)^opt\s*\{$`)
	zenParRe      = regexp.MustCompile(`(?i)^par\s*\{$`)
	zenNewRe      = regexp.MustCompile(`^(?:(\w+)\s*=\s*)?new\s+(\w+)\s*\(`)
	zenAsyncRe    = regexp.MustCompile(`^(\w+)\s*->\s*(\w+)\s*:\s*(.+)$`)
	zenSyncRe     = regexp.MustCompile(`^(?:(\w+)\s*=\s*)?(\w+)\.(\w+)\(`)
	zenSelfCallRe = regexp.MustCompile(`^(\w+)\(`)
	zenAtReturnRe = regexp.MustCompile(`(?i)^@return\s+(\w+)\s*->\s*(\w+)\s*:\s*(.+)$`)
	zenIdentRe    = regexp.MustCompile(`^\w+$`)
)

// zenControlKeywords lists keywords that should not be treated as self-calls.
//
//nolint:gochecknoglobals // package-level lookup table is idiomatic for constant keyword sets.
var zenControlKeywords = map[string]bool{
	"if": true, "else": true, "while": true, "for": true,
	"foreach": true, "loop": true, "try": true, "catch": true,
	"finally": true, "opt": true, "par": true, "break": true,
	"critical": true, "new": true, "return": true, "title": true,
	"group": true, "zenuml": true,
}

//nolint:gocognit,funlen,gocyclo,cyclop,maintidx // ZenUML parsing has inherent complexity from 20+ syntax patterns and block nesting.
func parseZenUML(input string) (*ParseOutput, error) {
	graph := ir.NewGraph()
	graph.Kind = ir.ZenUML

	lines := zenPreprocess(input)

	pIdx := map[string]int{}
	var stack []zenBlock
	caller := ""
	var currentBox *ir.SeqBox
	inGroup := false

	ensure := func(id string) {
		if _, ok := pIdx[id]; ok {
			return
		}
		pIdx[id] = len(graph.Participants)
		graph.Participants = append(graph.Participants, &ir.SeqParticipant{
			ID:   id,
			Kind: ir.ParticipantBox,
		})
		if inGroup && currentBox != nil {
			currentBox.Participants = append(currentBox.Participants, id)
		}
	}

	find := func(id string) *ir.SeqParticipant {
		if idx, ok := pIdx[id]; ok {
			return graph.Participants[idx]
		}
		return nil
	}

	emit := func(ev *ir.SeqEvent) {
		graph.Events = append(graph.Events, ev)
	}

	// closeBlock pops the top block from the stack and emits appropriate
	// close events. remainder is the text remaining after the '}' on the
	// same line, and nextLine is the following line (for split-line
	// continuations like "}\n else {"). Both are used to detect
	// continuations (else, catch, finally).
	closeBlock := func(remainder, nextLine string) {
		if len(stack) == 0 {
			return
		}
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Check both same-line remainder and next-line for continuations.
		isContinuation := func(keywords ...string) bool {
			for _, source := range []string{remainder, nextLine} {
				lower := strings.ToLower(strings.TrimSpace(source))
				for _, kw := range keywords {
					if strings.HasPrefix(lower, kw) {
						return true
					}
				}
			}
			return false
		}

		switch top.kind {
		case zenMsgBlock:
			emit(&ir.SeqEvent{Kind: ir.EvDeactivate, Target: top.target})
			caller = top.caller

		case zenIfBlock, zenElseIfBlock:
			if isContinuation("else") {
				return // continuation follows, don't close frame
			}
			emit(&ir.SeqEvent{Kind: ir.EvFrameEnd})

		case zenTryBlock, zenCatchBlock:
			if isContinuation("catch", "finally") {
				return // continuation follows
			}
			emit(&ir.SeqEvent{Kind: ir.EvFrameEnd})

		case zenGroupBlock:
			if currentBox != nil {
				graph.Boxes = append(graph.Boxes, currentBox)
				currentBox = nil
			}
			inGroup = false

		default: // zenElseBlock, zenFinallyBlock, zenLoopBlock, zenOptBlock, zenParBlock
			emit(&ir.SeqEvent{Kind: ir.EvFrameEnd})
		}
	}

	for lineIdx := range lines {
		line := lines[lineIdx]

		// Process leading close braces. Each '}' closes the top block.
		for strings.HasPrefix(line, "}") {
			line = strings.TrimSpace(line[1:])
			nextLine := ""
			if line == "" && lineIdx+1 < len(lines) {
				nextLine = lines[lineIdx+1]
			}
			closeBlock(line, nextLine)
		}
		if line == "" {
			continue
		}

		lower := strings.ToLower(line)

		// Skip header.
		if lower == "zenuml" {
			continue
		}

		// title (parsed but not stored — sequence layout doesn't use it).
		if strings.HasPrefix(lower, "title ") {
			continue
		}

		// @Starter(Participant)
		if match := zenStarterRe.FindStringSubmatch(line); match != nil {
			caller = match[1]
			ensure(caller)
			continue
		}

		// @return A->B: text
		if match := zenAtReturnRe.FindStringSubmatch(line); match != nil {
			from := match[1]
			to := match[2]
			text := strings.TrimSpace(match[3])
			ensure(from)
			ensure(to)
			emit(&ir.SeqEvent{
				Kind: ir.EvMessage,
				Message: &ir.SeqMessage{
					From: from,
					To:   to,
					Text: text,
					Kind: ir.MsgDottedArrow,
				},
			})
			continue
		}

		// Participant annotation: @Actor A, @Database DB as MyDB
		if match := zenAnnotRe.FindStringSubmatch(line); match != nil {
			kind := seqKindFromString(match[1])
			id := match[2]
			alias := strings.TrimSpace(match[3])
			ensure(id)
			participant := find(id)
			participant.Kind = kind
			if alias != "" {
				participant.Alias = alias
			}
			continue
		}

		// Alias: A as Alice
		if match := zenAliasRe.FindStringSubmatch(line); match != nil {
			id := match[1]
			alias := strings.TrimSpace(match[2])
			ensure(id)
			find(id).Alias = alias
			continue
		}

		// group Name {
		if match := zenGroupRe.FindStringSubmatch(line); match != nil {
			name := strings.TrimSpace(match[1])
			currentBox = &ir.SeqBox{Label: name}
			inGroup = true
			stack = append(stack, zenBlock{kind: zenGroupBlock, caller: caller})
			continue
		}

		// else if(cond) {
		if match := zenElseIfRe.FindStringSubmatch(line); match != nil {
			cond := match[1]
			emit(&ir.SeqEvent{
				Kind:  ir.EvFrameMiddle,
				Frame: &ir.SeqFrame{Kind: ir.FrameAlt, Label: cond},
			})
			stack = append(stack, zenBlock{kind: zenElseIfBlock, caller: caller})
			continue
		}

		// else {
		if (lower == "else {" || lower == "else{") && strings.HasSuffix(line, "{") {
			emit(&ir.SeqEvent{
				Kind:  ir.EvFrameMiddle,
				Frame: &ir.SeqFrame{Kind: ir.FrameAlt, Label: "else"},
			})
			stack = append(stack, zenBlock{kind: zenElseBlock, caller: caller})
			continue
		}

		// if(cond) {
		if match := zenIfRe.FindStringSubmatch(line); match != nil {
			cond := match[1]
			emit(&ir.SeqEvent{
				Kind:  ir.EvFrameStart,
				Frame: &ir.SeqFrame{Kind: ir.FrameAlt, Label: cond},
			})
			stack = append(stack, zenBlock{kind: zenIfBlock, caller: caller})
			continue
		}

		// while/for/forEach/loop
		if match := zenLoopRe.FindStringSubmatch(line); match != nil {
			label := match[2]
			emit(&ir.SeqEvent{
				Kind:  ir.EvFrameStart,
				Frame: &ir.SeqFrame{Kind: ir.FrameLoop, Label: label},
			})
			stack = append(stack, zenBlock{kind: zenLoopBlock, caller: caller})
			continue
		}

		// try { or try() {
		if zenTryRe.MatchString(line) {
			emit(&ir.SeqEvent{
				Kind:  ir.EvFrameStart,
				Frame: &ir.SeqFrame{Kind: ir.FrameAlt, Label: "try"},
			})
			stack = append(stack, zenBlock{kind: zenTryBlock, caller: caller})
			continue
		}

		// catch { or catch() {
		if zenCatchRe.MatchString(line) {
			emit(&ir.SeqEvent{
				Kind:  ir.EvFrameMiddle,
				Frame: &ir.SeqFrame{Kind: ir.FrameAlt, Label: "catch"},
			})
			stack = append(stack, zenBlock{kind: zenCatchBlock, caller: caller})
			continue
		}

		// finally { or finally() {
		if zenFinallyRe.MatchString(line) {
			emit(&ir.SeqEvent{
				Kind:  ir.EvFrameMiddle,
				Frame: &ir.SeqFrame{Kind: ir.FrameAlt, Label: "finally"},
			})
			stack = append(stack, zenBlock{kind: zenFinallyBlock, caller: caller})
			continue
		}

		// opt {
		if zenOptRe.MatchString(line) {
			emit(&ir.SeqEvent{
				Kind:  ir.EvFrameStart,
				Frame: &ir.SeqFrame{Kind: ir.FrameOpt},
			})
			stack = append(stack, zenBlock{kind: zenOptBlock, caller: caller})
			continue
		}

		// par {
		if zenParRe.MatchString(line) {
			emit(&ir.SeqEvent{
				Kind:  ir.EvFrameStart,
				Frame: &ir.SeqFrame{Kind: ir.FramePar},
			})
			stack = append(stack, zenBlock{kind: zenParBlock, caller: caller})
			continue
		}

		// return [value]
		if lower == "return" || strings.HasPrefix(lower, "return ") {
			retVal := strings.TrimSpace(line[len("return"):])
			// Walk stack to find enclosing message block for the return target.
			// If no message block exists (orphan return), skip silently — the
			// return has no source/destination in the sequence diagram.
			for stackIdx := len(stack) - 1; stackIdx >= 0; stackIdx-- {
				if stack[stackIdx].kind == zenMsgBlock {
					emit(&ir.SeqEvent{
						Kind: ir.EvMessage,
						Message: &ir.SeqMessage{
							From: stack[stackIdx].target,
							To:   stack[stackIdx].caller,
							Text: retVal,
							Kind: ir.MsgDottedArrow,
						},
					})
					break
				}
			}
			continue
		}

		// new Object() or obj = new Object()
		if match := zenNewRe.FindStringSubmatchIndex(line); match != nil {
			varName := ""
			if match[2] >= 0 {
				varName = line[match[2]:match[3]]
			}
			className := line[match[4]:match[5]]
			openIdx := match[1] - 1 // position of '('
			args, hasBlock, ok := zenParseCallArgs(line, openIdx)
			if !ok {
				continue
			}

			ensure(className)

			text := "new " + className + "(" + args + ")"
			if varName != "" {
				text = varName + " = " + text
			}

			emit(&ir.SeqEvent{Kind: ir.EvCreate, Target: className})
			if caller != "" {
				emit(&ir.SeqEvent{
					Kind: ir.EvMessage,
					Message: &ir.SeqMessage{
						From: caller,
						To:   className,
						Text: text,
						Kind: ir.MsgSolidArrow,
					},
				})
			}

			if hasBlock {
				emit(&ir.SeqEvent{Kind: ir.EvActivate, Target: className})
				stack = append(stack, zenBlock{kind: zenMsgBlock, caller: caller, target: className})
				caller = className
			}
			continue
		}

		// Async message: A->B: text
		if match := zenAsyncRe.FindStringSubmatch(line); match != nil {
			from := match[1]
			to := match[2]
			text := strings.TrimSpace(match[3])
			ensure(from)
			ensure(to)
			emit(&ir.SeqEvent{
				Kind: ir.EvMessage,
				Message: &ir.SeqMessage{
					From: from,
					To:   to,
					Text: text,
					Kind: ir.MsgSolidOpen,
				},
			})
			continue
		}

		// Sync message: A.method() or result = A.method() or A.method() {
		if match := zenSyncRe.FindStringSubmatchIndex(line); match != nil {
			varName := ""
			if match[2] >= 0 {
				varName = line[match[2]:match[3]]
			}
			target := line[match[4]:match[5]]
			methodName := line[match[6]:match[7]]
			openIdx := match[1] - 1 // position of '('
			args, hasBlock, ok := zenParseCallArgs(line, openIdx)
			if !ok {
				continue
			}

			from := caller
			if from == "" {
				from = target // self-call at top level
			}
			ensure(target)
			if from != target {
				ensure(from)
			}

			text := methodName + "(" + args + ")"
			if varName != "" {
				text = varName + " = " + text
			}

			emit(&ir.SeqEvent{
				Kind: ir.EvMessage,
				Message: &ir.SeqMessage{
					From: from,
					To:   target,
					Text: text,
					Kind: ir.MsgSolidArrow,
				},
			})

			if hasBlock {
				emit(&ir.SeqEvent{Kind: ir.EvActivate, Target: target})
				stack = append(stack, zenBlock{kind: zenMsgBlock, caller: caller, target: target})
				caller = target
			}
			continue
		}

		// Self-call: method() or method(args) or method() {
		if match := zenSelfCallRe.FindStringSubmatchIndex(line); match != nil && caller != "" {
			methodName := line[match[2]:match[3]]

			// Skip control keywords.
			if zenControlKeywords[strings.ToLower(methodName)] {
				continue
			}

			openIdx := match[1] - 1 // position of '('
			args, hasBlock, ok := zenParseCallArgs(line, openIdx)
			if !ok {
				continue
			}

			text := methodName + "(" + args + ")"

			emit(&ir.SeqEvent{
				Kind: ir.EvMessage,
				Message: &ir.SeqMessage{
					From: caller,
					To:   caller,
					Text: text,
					Kind: ir.MsgSolidArrow,
				},
			})

			if hasBlock {
				emit(&ir.SeqEvent{Kind: ir.EvActivate, Target: caller})
				stack = append(stack, zenBlock{kind: zenMsgBlock, caller: caller, target: caller})
			}
			continue
		}

		// Bare identifier = participant declaration.
		if zenIdentRe.MatchString(line) && !zenControlKeywords[lower] {
			ensure(line)
			continue
		}
	}

	// Validate: unclosed blocks are structural errors.
	if len(stack) > 0 {
		return nil, &ParseError{
			Diagram: "zenuml",
			Message: "unclosed block (missing \"}\")",
		}
	}

	return &ParseOutput{Graph: graph}, nil
}

// zenPreprocess strips // comments (quote-aware) and %% comments,
// filters empty lines, and returns cleaned lines.
func zenPreprocess(input string) []string {
	var lines []string
	for _, rawLine := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(rawLine)
		if trimmed == "" {
			continue
		}
		// Strip // line comments (quote-aware).
		trimmed = zenStripLineComment(trimmed)
		if trimmed == "" {
			continue
		}
		// Also strip %% comments.
		if strings.HasPrefix(trimmed, "%%") {
			continue
		}
		without := stripTrailingComment(trimmed)
		if without == "" {
			continue
		}
		lines = append(lines, without)
	}
	return lines
}

// zenFindBalancedParen finds the balanced closing ')' starting from the open
// paren at position openIdx. Returns the index of the closing ')' or -1.
func zenFindBalancedParen(line string, openIdx int) int {
	depth := 0
	for pos := openIdx; pos < len(line); pos++ {
		switch line[pos] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return pos
			}
		}
	}
	return -1
}

// zenParseCallArgs extracts the argument string and trailing content from a
// line starting at the open paren position. Returns (args, hasBlock) where
// hasBlock indicates a trailing '{'.
//
//nolint:nonamedreturns // named returns clarify the multi-value return.
func zenParseCallArgs(line string, openIdx int) (args string, hasBlock bool, ok bool) {
	closeIdx := zenFindBalancedParen(line, openIdx)
	if closeIdx < 0 {
		return "", false, false
	}
	args = line[openIdx+1 : closeIdx]
	rest := strings.TrimSpace(line[closeIdx+1:])
	hasBlock = rest == "{"
	return args, hasBlock, true
}

// zenStripLineComment removes // comments while respecting quoted strings.
func zenStripLineComment(line string) string {
	inQuote := false
	var quoteChar byte
	for pos := range len(line) {
		ch := line[pos]
		switch {
		case !inQuote && (ch == '"' || ch == '\''):
			inQuote = true
			quoteChar = ch
		case inQuote && ch == quoteChar:
			inQuote = false
		case !inQuote && ch == '/' && pos+1 < len(line) && line[pos+1] == '/':
			return strings.TrimSpace(line[:pos])
		}
	}
	return line
}
