package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseZenUMLBasicMessage(t *testing.T) {
	input := `zenuml
@Starter(Client)
Server.request() {
  Database.query()
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph
	if g.Kind != ir.ZenUML {
		t.Fatalf("kind = %v, want ZenUML", g.Kind)
	}
	if len(g.Participants) != 3 {
		t.Fatalf("participants = %d, want 3", len(g.Participants))
	}
	if g.Participants[0].ID != "Client" {
		t.Errorf("participant[0] = %q, want Client", g.Participants[0].ID)
	}
	if g.Participants[1].ID != "Server" {
		t.Errorf("participant[1] = %q, want Server", g.Participants[1].ID)
	}

	// Events: message(Client→Server), activate(Server),
	//         message(Server→Database), deactivate(Server)
	wantEvents := []ir.SeqEventKind{ir.EvMessage, ir.EvActivate, ir.EvMessage, ir.EvDeactivate}
	if len(g.Events) != len(wantEvents) {
		t.Fatalf("events = %d, want %d; got: %v", len(g.Events), len(wantEvents), eventKinds(g.Events))
	}
	for i, want := range wantEvents {
		if g.Events[i].Kind != want {
			t.Errorf("event[%d].Kind = %v, want %v", i, g.Events[i].Kind, want)
		}
	}

	// First message: Client → Server
	msg := g.Events[0].Message
	if msg.From != "Client" || msg.To != "Server" {
		t.Errorf("msg[0] = %s→%s, want Client→Server", msg.From, msg.To)
	}
	if msg.Text != "request()" {
		t.Errorf("msg[0].Text = %q, want %q", msg.Text, "request()")
	}
	if msg.Kind != ir.MsgSolidArrow {
		t.Errorf("msg[0].Kind = %v, want MsgSolidArrow", msg.Kind)
	}

	// Nested message: Server → Database (caller is Server inside the block)
	msg2 := g.Events[2].Message
	if msg2.From != "Server" || msg2.To != "Database" {
		t.Errorf("msg[2] = %s→%s, want Server→Database", msg2.From, msg2.To)
	}
}

func TestParseZenUMLParticipantAnnotations(t *testing.T) {
	input := `zenuml
@Actor Client
@Database UserDB
@Boundary API as PublicAPI
Client.login()
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph

	tests := []struct {
		id    string
		kind  ir.SeqParticipantKind
		alias string
	}{
		{"Client", ir.ActorStickFigure, ""},
		{"UserDB", ir.ParticipantDatabase, ""},
		{"API", ir.ParticipantBoundary, "PublicAPI"},
	}
	for _, tt := range tests {
		found := false
		for _, p := range g.Participants {
			if p.ID == tt.id {
				found = true
				if p.Kind != tt.kind {
					t.Errorf("participant %s kind = %v, want %v", tt.id, p.Kind, tt.kind)
				}
				if p.Alias != tt.alias {
					t.Errorf("participant %s alias = %q, want %q", tt.id, p.Alias, tt.alias)
				}
				break
			}
		}
		if !found {
			t.Errorf("participant %s not found", tt.id)
		}
	}
}

func TestParseZenUMLAlias(t *testing.T) {
	input := `zenuml
A as Alice
B as Bob
A.greet()
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if out.Graph.Participants[0].Alias != "Alice" {
		t.Errorf("alias = %q, want Alice", out.Graph.Participants[0].Alias)
	}
}

func TestParseZenUMLAsyncMessage(t *testing.T) {
	input := `zenuml
Alice->Bob: How are you?
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph
	if len(g.Events) != 1 {
		t.Fatalf("events = %d, want 1", len(g.Events))
	}
	msg := g.Events[0].Message
	if msg.From != "Alice" || msg.To != "Bob" {
		t.Errorf("msg = %s→%s, want Alice→Bob", msg.From, msg.To)
	}
	if msg.Kind != ir.MsgSolidOpen {
		t.Errorf("kind = %v, want MsgSolidOpen", msg.Kind)
	}
	if msg.Text != "How are you?" {
		t.Errorf("text = %q, want %q", msg.Text, "How are you?")
	}
}

func TestParseZenUMLReturn(t *testing.T) {
	input := `zenuml
@Starter(Client)
Server.process() {
  return success
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph

	// Events: message(Client→Server), activate(Server), message(Server→Client: success), deactivate(Server)
	if len(g.Events) != 4 {
		t.Fatalf("events = %d, want 4", len(g.Events))
	}

	// Return message: dotted arrow from Server back to Client.
	ret := g.Events[2].Message
	if ret.From != "Server" || ret.To != "Client" {
		t.Errorf("return = %s→%s, want Server→Client", ret.From, ret.To)
	}
	if ret.Kind != ir.MsgDottedArrow {
		t.Errorf("return kind = %v, want MsgDottedArrow", ret.Kind)
	}
	if ret.Text != "success" {
		t.Errorf("return text = %q, want %q", ret.Text, "success")
	}
}

func TestParseZenUMLAssignment(t *testing.T) {
	input := `zenuml
@Starter(Client)
result = Server.query(id)
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	msg := out.Graph.Events[0].Message
	if msg.Text != "result = query(id)" {
		t.Errorf("text = %q, want %q", msg.Text, "result = query(id)")
	}
}

func TestParseZenUMLIfElse(t *testing.T) {
	input := `zenuml
@Starter(Client)
if(authenticated) {
  Server.allow()
} else {
  Server.deny()
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph

	// Events: FrameStart(alt, "authenticated"), message, FrameMiddle("else"), message, FrameEnd
	wantKinds := []ir.SeqEventKind{
		ir.EvFrameStart,
		ir.EvMessage,
		ir.EvFrameMiddle,
		ir.EvMessage,
		ir.EvFrameEnd,
	}
	if len(g.Events) != len(wantKinds) {
		t.Fatalf("events = %d, want %d; got kinds: %v", len(g.Events), len(wantKinds), eventKinds(g.Events))
	}
	for i, want := range wantKinds {
		if g.Events[i].Kind != want {
			t.Errorf("event[%d] = %v, want %v", i, g.Events[i].Kind, want)
		}
	}

	// Frame start label.
	if g.Events[0].Frame.Label != "authenticated" {
		t.Errorf("frame label = %q, want %q", g.Events[0].Frame.Label, "authenticated")
	}
	if g.Events[0].Frame.Kind != ir.FrameAlt {
		t.Errorf("frame kind = %v, want FrameAlt", g.Events[0].Frame.Kind)
	}
}

func TestParseZenUMLElseIf(t *testing.T) {
	input := `zenuml
@Starter(A)
if(x > 0) {
  B.positive()
} else if(x == 0) {
  B.zero()
} else {
  B.negative()
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph

	wantKinds := []ir.SeqEventKind{
		ir.EvFrameStart,  // if(x > 0)
		ir.EvMessage,     // B.positive()
		ir.EvFrameMiddle, // else if(x == 0)
		ir.EvMessage,     // B.zero()
		ir.EvFrameMiddle, // else
		ir.EvMessage,     // B.negative()
		ir.EvFrameEnd,
	}
	if len(g.Events) != len(wantKinds) {
		t.Fatalf("events = %d, want %d; got: %v", len(g.Events), len(wantKinds), eventKinds(g.Events))
	}
	for i, want := range wantKinds {
		if g.Events[i].Kind != want {
			t.Errorf("event[%d] = %v, want %v", i, g.Events[i].Kind, want)
		}
	}

	// Check middle frame labels.
	if g.Events[2].Frame.Label != "x == 0" {
		t.Errorf("else-if label = %q, want %q", g.Events[2].Frame.Label, "x == 0")
	}
}

func TestParseZenUMLWhileLoop(t *testing.T) {
	input := `zenuml
@Starter(Client)
while(hasMore) {
  Server.fetch()
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph

	wantKinds := []ir.SeqEventKind{ir.EvFrameStart, ir.EvMessage, ir.EvFrameEnd}
	if len(g.Events) != len(wantKinds) {
		t.Fatalf("events = %d, want %d", len(g.Events), len(wantKinds))
	}
	if g.Events[0].Frame.Kind != ir.FrameLoop {
		t.Errorf("frame kind = %v, want FrameLoop", g.Events[0].Frame.Kind)
	}
	if g.Events[0].Frame.Label != "hasMore" {
		t.Errorf("frame label = %q, want %q", g.Events[0].Frame.Label, "hasMore")
	}
}

func TestParseZenUMLForLoop(t *testing.T) {
	input := `zenuml
@Starter(Client)
for(each item) {
  Server.process()
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if out.Graph.Events[0].Frame.Kind != ir.FrameLoop {
		t.Errorf("frame kind = %v, want FrameLoop", out.Graph.Events[0].Frame.Kind)
	}
}

func TestParseZenUMLTryCatchFinally(t *testing.T) {
	input := `zenuml
@Starter(Client)
try {
  Server.riskyOp()
} catch {
  Logger.error()
} finally {
  Server.cleanup()
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph

	wantKinds := []ir.SeqEventKind{
		ir.EvFrameStart,  // try
		ir.EvMessage,     // Server.riskyOp()
		ir.EvFrameMiddle, // catch
		ir.EvMessage,     // Logger.error()
		ir.EvFrameMiddle, // finally
		ir.EvMessage,     // Server.cleanup()
		ir.EvFrameEnd,
	}
	if len(g.Events) != len(wantKinds) {
		t.Fatalf("events = %d, want %d; got: %v", len(g.Events), len(wantKinds), eventKinds(g.Events))
	}
	for i, want := range wantKinds {
		if g.Events[i].Kind != want {
			t.Errorf("event[%d] = %v, want %v", i, g.Events[i].Kind, want)
		}
	}

	if g.Events[0].Frame.Label != "try" {
		t.Errorf("try label = %q, want %q", g.Events[0].Frame.Label, "try")
	}
	if g.Events[2].Frame.Label != "catch" {
		t.Errorf("catch label = %q, want %q", g.Events[2].Frame.Label, "catch")
	}
	if g.Events[4].Frame.Label != "finally" {
		t.Errorf("finally label = %q, want %q", g.Events[4].Frame.Label, "finally")
	}
}

func TestParseZenUMLTryParenSyntax(t *testing.T) {
	input := `zenuml
@Starter(Client)
try() {
  Server.op()
} catch() {
  Logger.log()
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph
	if g.Events[0].Kind != ir.EvFrameStart {
		t.Errorf("event[0] = %v, want EvFrameStart", g.Events[0].Kind)
	}
}

func TestParseZenUMLOptPar(t *testing.T) {
	input := `zenuml
@Starter(Client)
opt {
  Server.optional()
}
par {
  A.task1()
  B.task2()
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph

	// opt: FrameStart(Opt), message, FrameEnd, par: FrameStart(Par), message, message, FrameEnd
	if len(g.Events) < 7 {
		t.Fatalf("events = %d, want >= 7", len(g.Events))
	}
	if g.Events[0].Frame.Kind != ir.FrameOpt {
		t.Errorf("event[0] frame = %v, want FrameOpt", g.Events[0].Frame.Kind)
	}
	if g.Events[3].Frame.Kind != ir.FramePar {
		t.Errorf("event[3] frame = %v, want FramePar", g.Events[3].Frame.Kind)
	}
}

func TestParseZenUMLObjectCreation(t *testing.T) {
	input := `zenuml
@Starter(Client)
receipt = new Receipt()
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph

	// Events: EvCreate(Receipt), EvMessage(Client→Receipt)
	if len(g.Events) != 2 {
		t.Fatalf("events = %d, want 2; got: %v", len(g.Events), eventKinds(g.Events))
	}
	if g.Events[0].Kind != ir.EvCreate {
		t.Errorf("event[0] = %v, want EvCreate", g.Events[0].Kind)
	}
	if g.Events[0].Target != "Receipt" {
		t.Errorf("create target = %q, want Receipt", g.Events[0].Target)
	}
	msg := g.Events[1].Message
	if msg.Text != "receipt = new Receipt()" {
		t.Errorf("text = %q, want %q", msg.Text, "receipt = new Receipt()")
	}
}

func TestParseZenUMLSelfCall(t *testing.T) {
	input := `zenuml
@Starter(Client)
Server.handle() {
  validate()
  process()
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph

	// Events: msg(Client→Server), activate(Server), msg(Server→Server: validate()),
	//         msg(Server→Server: process()), deactivate(Server)
	wantKinds := []ir.SeqEventKind{
		ir.EvMessage, ir.EvActivate, ir.EvMessage, ir.EvMessage, ir.EvDeactivate,
	}
	if len(g.Events) != len(wantKinds) {
		t.Fatalf("events = %d, want %d; got: %v", len(g.Events), len(wantKinds), eventKinds(g.Events))
	}
	for i, want := range wantKinds {
		if g.Events[i].Kind != want {
			t.Errorf("event[%d] = %v, want %v", i, g.Events[i].Kind, want)
		}
	}

	// Self-calls should be from Server to Server.
	for _, idx := range []int{2, 3} {
		msg := g.Events[idx].Message
		if msg.From != "Server" || msg.To != "Server" {
			t.Errorf("event[%d] = %s→%s, want Server→Server", idx, msg.From, msg.To)
		}
	}
}

func TestParseZenUMLGroup(t *testing.T) {
	input := `zenuml
group BusinessService {
  @Actor A
  @Database B
}
A.call()
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph
	if len(g.Boxes) != 1 {
		t.Fatalf("boxes = %d, want 1", len(g.Boxes))
	}
	if g.Boxes[0].Label != "BusinessService" {
		t.Errorf("box label = %q, want BusinessService", g.Boxes[0].Label)
	}
	if len(g.Boxes[0].Participants) != 2 {
		t.Errorf("box participants = %d, want 2", len(g.Boxes[0].Participants))
	}
}

func TestParseZenUMLAtReturn(t *testing.T) {
	input := `zenuml
@Starter(Client)
Server.process()
@return Server->Client: done
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph
	if len(g.Events) != 2 {
		t.Fatalf("events = %d, want 2", len(g.Events))
	}
	ret := g.Events[1].Message
	if ret.From != "Server" || ret.To != "Client" {
		t.Errorf("@return = %s→%s, want Server→Client", ret.From, ret.To)
	}
	if ret.Kind != ir.MsgDottedArrow {
		t.Errorf("kind = %v, want MsgDottedArrow", ret.Kind)
	}
}

func TestParseZenUMLComments(t *testing.T) {
	input := `zenuml
// This is a comment
@Starter(Client)
Server.process() // inline comment
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Graph.Events) != 1 {
		t.Fatalf("events = %d, want 1", len(out.Graph.Events))
	}
}

func TestParseZenUMLNestedBlocks(t *testing.T) {
	input := `zenuml
@Starter(Client)
OrderController.post(payload) {
  OrderService.create(payload) {
    order = new Order()
    if(order != null) {
      par {
        PurchaseService.createPO()
        InvoiceService.createInvoice()
      }
    }
  }
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph

	// Verify participants were created.
	names := make(map[string]bool)
	for _, p := range g.Participants {
		names[p.ID] = true
	}
	for _, want := range []string{"Client", "OrderController", "OrderService", "Order", "PurchaseService", "InvoiceService"} {
		if !names[want] {
			t.Errorf("missing participant %q", want)
		}
	}

	// Verify frame events exist.
	var frameStarts, frameEnds int
	for _, ev := range g.Events {
		if ev.Kind == ir.EvFrameStart {
			frameStarts++
		}
		if ev.Kind == ir.EvFrameEnd {
			frameEnds++
		}
	}
	if frameStarts != 2 { // if + par
		t.Errorf("frame starts = %d, want 2", frameStarts)
	}
	if frameEnds != 2 {
		t.Errorf("frame ends = %d, want 2", frameEnds)
	}
}

func TestParseZenUMLNoStarter(t *testing.T) {
	// Without @Starter, first A.method() should be a self-call.
	input := `zenuml
A.process()
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph
	if len(g.Events) != 1 {
		t.Fatalf("events = %d, want 1", len(g.Events))
	}
	msg := g.Events[0].Message
	if msg.From != "A" || msg.To != "A" {
		t.Errorf("msg = %s→%s, want A→A (self-call)", msg.From, msg.To)
	}
}

func TestParseZenUMLBareLoop(t *testing.T) {
	input := `zenuml
@Starter(A)
loop {
  B.fetch()
}
`
	out, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	g := out.Graph
	if g.Events[0].Kind != ir.EvFrameStart {
		t.Errorf("event[0] = %v, want EvFrameStart", g.Events[0].Kind)
	}
	if g.Events[0].Frame.Kind != ir.FrameLoop {
		t.Errorf("frame kind = %v, want FrameLoop", g.Events[0].Frame.Kind)
	}
}

func TestZenPreprocess(t *testing.T) {
	input := `zenuml
// full line comment
@Starter(Client)
Server.call() // inline comment
%% mermaid comment
`
	lines := zenPreprocess(input)
	want := []string{"zenuml", "@Starter(Client)", "Server.call()"}
	if len(lines) != len(want) {
		t.Fatalf("lines = %d, want %d; got: %v", len(lines), len(want), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line[%d] = %q, want %q", i, lines[i], w)
		}
	}
}

func TestZenStripLineComment(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"code // comment", "code"},
		{"no comment", "no comment"},
		{`"has // in string" // real`, `"has // in string"`},
		{"// full line", ""},
		{`'str // not' // yes`, `'str // not'`},
	}
	for _, tt := range tests {
		got := zenStripLineComment(tt.input)
		if got != tt.want {
			t.Errorf("zenStripLineComment(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// eventKinds is a test helper that returns event kinds for error messages.
func eventKinds(events []*ir.SeqEvent) []ir.SeqEventKind {
	kinds := make([]ir.SeqEventKind, len(events))
	for i, ev := range events {
		kinds[i] = ev.Kind
	}
	return kinds
}
