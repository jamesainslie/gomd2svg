package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

// Class diagram relationship regex patterns compiled at package level.
var (
	// classRelAllRe handles all 8 relationship types.
	classRelAllRe = regexp.MustCompile(
		`^(\S+)\s+` + // left
			`(?:"([^"]*?)"\s+)?` + // optional start cardinality
			`(<\|--|\*--|o--|-->|\.\.>|\.\.\|>|--|\.\.)` + // arrow
			`\s+(?:"([^"]*?)"\s+)?` + // optional end cardinality
			`(\S+)` + // right
			`(?:\s*:\s*(.+))?$`, // optional label
	)

	// classAnnotationRe matches <<annotation>> ClassName
	classAnnotationRe = regexp.MustCompile(`^<<(\w+)>>\s+(\w+)$`)

	// classBodyDeclRe matches "class ClassName {" or "class ClassName"
	classBodyDeclRe = regexp.MustCompile(`^class\s+(\w+(?:~[^~]+~)?)\s*(\{?)$`)

	// classColonMemberRe matches "ClassName : memberDef"
	classColonMemberRe = regexp.MustCompile(`^(\w+)\s*:\s*(.+)$`)

	// classNoteForRe matches "note for ClassName "text""
	classNoteForRe = regexp.MustCompile(`^note\s+for\s+(\w+)\s+"([^"]+)"$`)

	// classNoteRe matches "note "text""
	classNoteRe = regexp.MustCompile(`^note\s+"([^"]+)"$`)

	// classNamespaceRe matches "namespace Name {"
	classNamespaceRe = regexp.MustCompile(`^namespace\s+(\w+)\s*\{$`)

	// classGenericRe extracts class name and generic param from "ClassName~T~"
	classGenericRe = regexp.MustCompile(`^(\w+)~([^~]+)~$`)

	// classDirectiveSkipRe matches lines to skip (classDef, style, etc.)
	classDirectiveSkipRe = regexp.MustCompile(`(?i)^(classdef|style|cssclass|click|callback|link)\b`)
)

// parseClass parses a Mermaid class diagram.
func parseClass(input string) (*ParseOutput, error) {
	graph := ir.NewGraph()
	graph.Kind = ir.Class

	lines := preprocessInput(input)

	var currentClass string
	braceDepth := 0
	inNamespace := false
	var currentNamespace *ir.Namespace
	namespaceBraceDepth := 0

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Skip header line.
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "classdiagram") {
			continue
		}

		// Skip directives.
		if classDirectiveSkipRe.MatchString(line) {
			continue
		}

		// Handle direction.
		if dir, ok := parseDirectionLine(line); ok {
			graph.Direction = dir
			continue
		}

		// Handle brace closing.
		if line == "}" {
			if braceDepth > 0 {
				braceDepth--
				if braceDepth == 0 {
					currentClass = ""
				}
				continue
			}
			if inNamespace {
				namespaceBraceDepth--
				if namespaceBraceDepth == 0 {
					inNamespace = false
					currentNamespace = nil
				}
				continue
			}
			continue
		}

		// Inside a class body: parse members.
		if braceDepth > 0 && currentClass != "" {
			member := parseClassMember(line)
			ensureClassMembers(graph, currentClass)
			if member.IsMethod {
				graph.Members[currentClass].Methods = append(graph.Members[currentClass].Methods, member)
			} else {
				graph.Members[currentClass].Attributes = append(graph.Members[currentClass].Attributes, member)
			}
			continue
		}

		// Namespace declaration.
		if caps := classNamespaceRe.FindStringSubmatch(line); caps != nil {
			ns := &ir.Namespace{Name: caps[1]}
			graph.Namespaces = append(graph.Namespaces, ns)
			currentNamespace = ns
			inNamespace = true
			namespaceBraceDepth = 1
			continue
		}

		// Inside namespace: handle class declarations.
		if inNamespace && currentNamespace != nil {
			if strings.HasPrefix(line, "class ") {
				className := extractClassName(strings.TrimPrefix(line, "class "))
				currentNamespace.Classes = append(currentNamespace.Classes, className)
				graph.EnsureNode(className, nil, nil)
				// Check if class has an opening brace on same line.
				trimmed := strings.TrimSpace(strings.TrimPrefix(line, "class "))
				if strings.HasSuffix(trimmed, "{") {
					braceDepth = 1
					currentClass = className
				}
				continue
			}
		}

		// Annotation: <<keyword>> ClassName
		if caps := classAnnotationRe.FindStringSubmatch(line); caps != nil {
			graph.Annotations[caps[2]] = caps[1]
			graph.EnsureNode(caps[2], nil, nil)
			continue
		}

		// Note for class.
		if caps := classNoteForRe.FindStringSubmatch(line); caps != nil {
			graph.Notes = append(graph.Notes, &ir.DiagramNote{
				Text:   caps[2],
				Target: caps[1],
			})
			continue
		}

		// Note standalone.
		if caps := classNoteRe.FindStringSubmatch(line); caps != nil {
			graph.Notes = append(graph.Notes, &ir.DiagramNote{
				Text: caps[1],
			})
			continue
		}

		// Class body declaration: class ClassName { or class ClassName
		if caps := classBodyDeclRe.FindStringSubmatch(line); caps != nil {
			className := extractClassName(caps[1])
			graph.EnsureNode(className, nil, nil)
			ensureClassMembers(graph, className)
			if caps[2] == "{" {
				braceDepth = 1
				currentClass = className
			}
			continue
		}

		// Relationship line.
		if caps := classRelAllRe.FindStringSubmatch(line); caps != nil {
			leftID := caps[1]
			startCard := caps[2]
			arrow := caps[3]
			endCard := caps[4]
			rightID := caps[5]
			edgeLabel := strings.TrimSpace(caps[6])

			graph.EnsureNode(leftID, nil, nil)
			graph.EnsureNode(rightID, nil, nil)

			edge := buildClassEdge(leftID, rightID, arrow)

			if startCard != "" {
				edge.StartLabel = &startCard
			}
			if endCard != "" {
				edge.EndLabel = &endCard
			}
			if edgeLabel != "" {
				edge.Label = &edgeLabel
			}

			graph.Edges = append(graph.Edges, edge)
			continue
		}

		// Colon member syntax: ClassName : memberDef
		if caps := classColonMemberRe.FindStringSubmatch(line); caps != nil {
			className := caps[1]
			memberDef := strings.TrimSpace(caps[2])
			graph.EnsureNode(className, nil, nil)
			ensureClassMembers(graph, className)
			member := parseClassMember(memberDef)
			if member.IsMethod {
				graph.Members[className].Methods = append(graph.Members[className].Methods, member)
			} else {
				graph.Members[className].Attributes = append(graph.Members[className].Attributes, member)
			}
			continue
		}
	}

	return &ParseOutput{Graph: graph}, nil
}

// extractClassName strips generic suffix ~T~ and returns the bare class name.
func extractClassName(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if caps := classGenericRe.FindStringSubmatch(trimmed); caps != nil {
		return caps[1]
	}
	// Remove trailing semicolons.
	return strings.TrimRight(trimmed, ";")
}

// ensureClassMembers initializes the Members entry for a class if missing.
func ensureClassMembers(graph *ir.Graph, className string) {
	if graph.Members[className] == nil {
		graph.Members[className] = &ir.ClassMembers{}
	}
}

// parseClassMember parses a single class member line (attribute or method).
// Handles visibility prefixes (+,-,#,~), method detection (parentheses),
// classifier suffixes (* for abstract, $ for static), and type/name extraction.
func parseClassMember(line string) ir.ClassMember {
	trimmed := strings.TrimSpace(line)
	member := ir.ClassMember{}

	// Check visibility prefix.
	if len(trimmed) > 0 {
		switch trimmed[0] {
		case '+':
			member.Visibility = ir.VisPublic
			trimmed = strings.TrimSpace(trimmed[1:])
		case '-':
			member.Visibility = ir.VisPrivate
			trimmed = strings.TrimSpace(trimmed[1:])
		case '#':
			member.Visibility = ir.VisProtected
			trimmed = strings.TrimSpace(trimmed[1:])
		case '~':
			member.Visibility = ir.VisPackage
			trimmed = strings.TrimSpace(trimmed[1:])
		}
	}

	// Check if method (contains parentheses).
	parenOpen := strings.Index(trimmed, "(")
	parenClose := strings.LastIndex(trimmed, ")")
	if parenOpen >= 0 && parenClose > parenOpen {
		member.IsMethod = true

		// Extract method name (everything before the parenthesis).
		namePart := strings.TrimSpace(trimmed[:parenOpen])
		member.Name = namePart
		member.Params = trimmed[parenOpen+1 : parenClose]

		// Everything after the closing paren.
		after := strings.TrimSpace(trimmed[parenClose+1:])

		// Check classifier suffix on method: )* or )$
		if len(after) > 0 && (after[0] == '*' || after[0] == '$') {
			if after[0] == '*' {
				member.Classifier = ir.ClassifierAbstract
			} else {
				member.Classifier = ir.ClassifierStatic
			}
			after = strings.TrimSpace(after[1:])
		}

		// Remaining is return type.
		if after != "" {
			member.Type = after
		}
	} else {
		// Attribute: [type] name[classifier]
		parts := strings.Fields(trimmed)
		if len(parts) >= 2 {
			member.Type = parts[0]
			nameStr := parts[1]

			// Check classifier suffix on attribute.
			if strings.HasSuffix(nameStr, "*") {
				member.Classifier = ir.ClassifierAbstract
				nameStr = strings.TrimSuffix(nameStr, "*")
			} else if strings.HasSuffix(nameStr, "$") {
				member.Classifier = ir.ClassifierStatic
				nameStr = strings.TrimSuffix(nameStr, "$")
			}
			member.Name = nameStr
		} else if len(parts) == 1 {
			nameStr := parts[0]
			if strings.HasSuffix(nameStr, "*") {
				member.Classifier = ir.ClassifierAbstract
				nameStr = strings.TrimSuffix(nameStr, "*")
			} else if strings.HasSuffix(nameStr, "$") {
				member.Classifier = ir.ClassifierStatic
				nameStr = strings.TrimSuffix(nameStr, "$")
			}
			member.Name = nameStr
		}
	}

	return member
}

// buildClassEdge creates an Edge for a class diagram relationship arrow.
func buildClassEdge(from, to, arrow string) *ir.Edge {
	edge := &ir.Edge{
		From: from,
		To:   to,
	}

	switch arrow {
	case "<|--": // inheritance (to inherits from)
		edge.Directed = true
		edge.ArrowStart = true
		edge.ArrowEnd = false
		kind := ir.ClosedTriangle
		edge.ArrowStartKind = &kind
		edge.Style = ir.Solid
	case "*--": // composition
		edge.Directed = true
		edge.ArrowStart = true
		edge.ArrowEnd = false
		kind := ir.FilledDiamond
		edge.ArrowStartKind = &kind
		edge.Style = ir.Solid
	case "o--": // aggregation
		edge.Directed = true
		edge.ArrowStart = true
		edge.ArrowEnd = false
		kind := ir.OpenDiamond
		edge.ArrowStartKind = &kind
		edge.Style = ir.Solid
	case "-->": // association
		edge.Directed = true
		edge.ArrowStart = false
		edge.ArrowEnd = true
		kind := ir.OpenTriangle
		edge.ArrowEndKind = &kind
		edge.Style = ir.Solid
	case "..>": // dependency
		edge.Directed = true
		edge.ArrowStart = false
		edge.ArrowEnd = true
		kind := ir.ClassDependency
		edge.ArrowEndKind = &kind
		edge.Style = ir.Dotted
	case "..|>": // realization
		edge.Directed = true
		edge.ArrowStart = false
		edge.ArrowEnd = true
		kind := ir.ClosedTriangle
		edge.ArrowEndKind = &kind
		edge.Style = ir.Dotted
	case "--": // solid link (undirected)
		edge.Directed = false
		edge.Style = ir.Solid
	case "..": // dashed link (undirected)
		edge.Directed = false
		edge.Style = ir.Dotted
	}

	return edge
}
