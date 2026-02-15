# Phase 11: ZenUML Design

**Goal:** Add ZenUML parser that produces the same sequence diagram IR, reusing existing layout and renderer.

**Architecture:** ZenUML is a code-like alternative to `sequenceDiagram` syntax. It uses `A.method()` instead of `A->>B: method` and curly braces for nesting. The parser converts ZenUML into the same `SeqParticipant`, `SeqEvent`, `SeqMessage`, `SeqFrame` IR types already used by the sequence diagram. Layout (`computeSequenceLayout`) and rendering (`renderSequence`) are reused without modification.

## Syntax Supported

### Participants
- Implicit: appear on first use
- Explicit: bare identifier on its own line (`Alice`)
- Aliases: `A as Alice`
- Annotations: `@Actor A`, `@Database DB`, `@Boundary B`, `@Control C`, `@Entity E`, `@Collections C`, `@Queue Q`
- Groups: `group Name { @Actor A; @Database B }`

### Messages
- Sync: `A.method()`, `A.method(args)`, `result = A.method()`
- Sync with block: `A.method() { ... }` — activates target, nested calls use target as caller
- Async: `A->B: message text`
- Self-call: bare `method()` inside a block — calls current participant
- Creation: `new Object()`, `obj = new Object(params)`
- Return: `return value` — dotted arrow back to caller

### Control Flow
- `if(cond) { } else if(cond) { } else { }` → FrameAlt with EvFrameMiddle dividers
- `while(cond) { }`, `for(cond) { }`, `forEach(cond) { }`, `loop { }` → FrameLoop
- `try { } catch { } finally { }` → FrameAlt with EvFrameMiddle dividers
- `opt { }` → FrameOpt
- `par { }` → FramePar

### Other
- `title Text` → graph title (not used in current layout, stored for future)
- `@Starter(Participant)` → sets initial caller
- `// comment` → stripped during preprocessing

## IR Mapping

| ZenUML Syntax | IR Event |
|---|---|
| `A.method()` | EvMessage(MsgSolidArrow, caller→A, "method()") |
| `A.method() {` | EvMessage + EvActivate(A), push caller stack |
| `}` closing message block | EvDeactivate(A), pop caller stack |
| `A->B: text` | EvMessage(MsgSolidOpen, A→B, "text") |
| `return val` | EvMessage(MsgDottedArrow, current→enclosing_caller, "val") |
| `new Obj()` | EvCreate(Obj) + EvMessage(MsgSolidArrow, caller→Obj) |
| `if(cond) {` | EvFrameStart(FrameAlt, "cond") |
| `} else {` | EvFrameMiddle(FrameAlt, "else") |
| `}` closing if/else | EvFrameEnd |
| `while(cond) {` | EvFrameStart(FrameLoop, "cond") |
| `try {` | EvFrameStart(FrameAlt, "try") |
| `} catch {` | EvFrameMiddle(FrameAlt, "catch") |
| `} finally {` | EvFrameMiddle(FrameAlt, "finally") |
| `opt {` | EvFrameStart(FrameOpt) |
| `par {` | EvFrameStart(FramePar) |
| `group Name {` | SeqBox{Label:"Name"} |

## Parser Design

Recursive line-based parser with a block stack:

1. **Preprocess**: Strip `//` comments (not `%%`), filter empties
2. **Block stack**: Each `{` pushes a block, each `}` pops. Blocks track kind (message/if/loop/try/etc.) and caller context
3. **Continuation handling**: `} else {`, `} catch {`, `} finally {` emit EvFrameMiddle instead of EvFrameEnd
4. **Caller tracking**: Nested `A.method() { B.call() }` sets caller to A inside the block, restores on close

## Dispatch Wiring

- `parser/parser.go`: `case ir.ZenUML: return parseZenUML(input)`
- `layout/layout.go`: `case ir.ZenUML: return computeSequenceLayout(g, th, cfg)` (full reuse)
- `render/svg.go`: No changes needed — dispatches on `layout.SequenceData` which ZenUML produces

## Files

| File | Action |
|---|---|
| `parser/zenuml.go` | Create — full ZenUML parser |
| `parser/zenuml_test.go` | Create — table-driven tests |
| `parser/parser.go` | Modify — add dispatch case |
| `layout/layout.go` | Modify — add dispatch case |
| `mermaid_test.go` | Modify — add integration tests |
| `testdata/fixtures/zenuml-*.mmd` | Create — test fixtures |
