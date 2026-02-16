package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
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
	zenNewRe      = regexp.MustCompile(`^(?:(\w+)\s*=\s*)?new\s+(\w+)\s*\([^)]*\)\s*(\{)?$`)
	zenAsyncRe    = regexp.MustCompile(`^(\w+)\s*->\s*(\w+)\s*:\s*(.+)$`)
	zenSyncRe     = regexp.MustCompile(`^(?:(\w+)\s*=\s*)?(\w+)\.(\w+)\(([^)]*)\)\s*(\{)?$`)
	zenSelfCallRe = regexp.MustCompile(`^(\w+)\(([^)]*)\)\s*(\{)?$`)
	zenAtReturnRe = regexp.MustCompile(`(?i)^@return\s+(\w+)\s*->\s*(\w+)\s*:\s*(.+)$`)
	zenIdentRe    = regexp.MustCompile(`^\w+$`)
)

// zenControlKeywords lists keywords that should not be treated as self-calls.
var zenControlKeywords = map[string]bool{
	"if": true, "else": true, "while": true, "for": true,
	"foreach": true, "loop": true, "try": true, "catch": true,
	"finally": true, "opt": true, "par": true, "break": true,
	"critical": true, "new": true, "return": true, "title": true,
	"group": true, "zenuml": true,
}

func parseZenUML(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.ZenUML

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
		pIdx[id] = len(g.Participants)
		g.Participants = append(g.Participants, &ir.SeqParticipant{
			ID:   id,
			Kind: ir.ParticipantBox,
		})
		if inGroup && currentBox != nil {
			currentBox.Participants = append(currentBox.Participants, id)
		}
	}

	find := func(id string) *ir.SeqParticipant {
		if idx, ok := pIdx[id]; ok {
			return g.Participants[idx]
		}
		return nil
	}

	emit := func(ev *ir.SeqEvent) {
		g.Events = append(g.Events, ev)
	}

	// closeBlock pops the top block from the stack and emits appropriate
	// close events. remainder is the text remaining after the '}' on the
	// same line, used to detect continuations (else, catch, finally).
	closeBlock := func(remainder string) {
		if len(stack) == 0 {
			return
		}
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch top.kind {
		case zenMsgBlock:
			emit(&ir.SeqEvent{Kind: ir.EvDeactivate, Target: top.target})
			caller = top.caller

		case zenIfBlock, zenElseIfBlock:
			lower := strings.ToLower(strings.TrimSpace(remainder))
			if strings.HasPrefix(lower, "else") {
				return // continuation follows, don't close frame
			}
			emit(&ir.SeqEvent{Kind: ir.EvFrameEnd})

		case zenTryBlock, zenCatchBlock:
			lower := strings.ToLower(strings.TrimSpace(remainder))
			if strings.HasPrefix(lower, "catch") || strings.HasPrefix(lower, "finally") {
				return // continuation follows
			}
			emit(&ir.SeqEvent{Kind: ir.EvFrameEnd})

		case zenGroupBlock:
			if currentBox != nil {
				g.Boxes = append(g.Boxes, currentBox)
				currentBox = nil
			}
			inGroup = false

		default: // zenElseBlock, zenFinallyBlock, zenLoopBlock, zenOptBlock, zenParBlock
			emit(&ir.SeqEvent{Kind: ir.EvFrameEnd})
		}
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Process leading close braces. Each '}' closes the top block.
		for strings.HasPrefix(line, "}") {
			line = strings.TrimSpace(line[1:])
			closeBlock(line)
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
		if m := zenStarterRe.FindStringSubmatch(line); m != nil {
			caller = m[1]
			ensure(caller)
			continue
		}

		// @return A->B: text
		if m := zenAtReturnRe.FindStringSubmatch(line); m != nil {
			from := m[1]
			to := m[2]
			text := strings.TrimSpace(m[3])
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
		if m := zenAnnotRe.FindStringSubmatch(line); m != nil {
			kind := seqKindFromString(m[1])
			id := m[2]
			alias := strings.TrimSpace(m[3])
			ensure(id)
			p := find(id)
			p.Kind = kind
			if alias != "" {
				p.Alias = alias
			}
			continue
		}

		// Alias: A as Alice
		if m := zenAliasRe.FindStringSubmatch(line); m != nil {
			id := m[1]
			alias := strings.TrimSpace(m[2])
			ensure(id)
			find(id).Alias = alias
			continue
		}

		// group Name {
		if m := zenGroupRe.FindStringSubmatch(line); m != nil {
			name := strings.TrimSpace(m[1])
			currentBox = &ir.SeqBox{Label: name}
			inGroup = true
			stack = append(stack, zenBlock{kind: zenGroupBlock, caller: caller})
			continue
		}

		// else if(cond) {
		if m := zenElseIfRe.FindStringSubmatch(line); m != nil {
			cond := m[1]
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
		if m := zenIfRe.FindStringSubmatch(line); m != nil {
			cond := m[1]
			emit(&ir.SeqEvent{
				Kind:  ir.EvFrameStart,
				Frame: &ir.SeqFrame{Kind: ir.FrameAlt, Label: cond},
			})
			stack = append(stack, zenBlock{kind: zenIfBlock, caller: caller})
			continue
		}

		// while/for/forEach/loop
		if m := zenLoopRe.FindStringSubmatch(line); m != nil {
			label := m[2]
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
			for j := len(stack) - 1; j >= 0; j-- {
				if stack[j].kind == zenMsgBlock {
					emit(&ir.SeqEvent{
						Kind: ir.EvMessage,
						Message: &ir.SeqMessage{
							From: stack[j].target,
							To:   stack[j].caller,
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
		if m := zenNewRe.FindStringSubmatch(line); m != nil {
			varName := m[1]
			className := m[2]
			hasBlock := m[3] == "{"

			ensure(className)

			text := "new " + className + "()"
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
		if m := zenAsyncRe.FindStringSubmatch(line); m != nil {
			from := m[1]
			to := m[2]
			text := strings.TrimSpace(m[3])
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
		if m := zenSyncRe.FindStringSubmatch(line); m != nil {
			varName := m[1]
			target := m[2]
			methodName := m[3]
			args := m[4]
			hasBlock := m[5] == "{"

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
		if m := zenSelfCallRe.FindStringSubmatch(line); m != nil && caller != "" {
			methodName := m[1]
			args := m[2]
			hasBlock := m[3] == "{"

			// Skip control keywords.
			if zenControlKeywords[strings.ToLower(methodName)] {
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

	// Close any unclosed blocks.
	for len(stack) > 0 {
		closeBlock("")
	}

	return &ParseOutput{Graph: g}, nil
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

// zenStripLineComment removes // comments while respecting quoted strings.
func zenStripLineComment(line string) string {
	inQuote := false
	var quoteChar byte
	for i := 0; i < len(line); i++ {
		ch := line[i]
		if !inQuote && (ch == '"' || ch == '\'') {
			inQuote = true
			quoteChar = ch
		} else if inQuote && ch == quoteChar {
			inQuote = false
		} else if !inQuote && ch == '/' && i+1 < len(line) && line[i+1] == '/' {
			return strings.TrimSpace(line[:i])
		}
	}
	return line
}
