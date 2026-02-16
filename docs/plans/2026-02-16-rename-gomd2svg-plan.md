# Rename gomd2svg to gomd2svg — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rename project from gomd2svg to gomd2svg across all code, docs, and infrastructure.

**Architecture:** Mechanical find-and-replace in dependency order: module path first, then package names, then directory renames, then docs.

**Tech Stack:** Go, git, gh CLI

---

### Task 1: Rename GitHub repo

**Step 1:** Rename via gh CLI

```bash
gh repo rename gomd2svg
```

**Step 2:** Update local git remote

```bash
git remote set-url origin https://github.com/jamesainslie/gomd2svg.git
```

**Step 3:** Verify

```bash
git remote get-url origin
gh repo view --json name -q .name
```

---

### Task 2: Update go.mod module path

**Files:**
- Modify: `go.mod`

**Step 1:** Change module declaration

```
module github.com/jamesainslie/gomd2svg
```

---

### Task 3: Replace all import paths

**Files:**
- Modify: all 166 files containing `github.com/jamesainslie/gomd2svg`

**Step 1:** Bulk find-and-replace

Replace `github.com/jamesainslie/gomd2svg` with `github.com/jamesainslie/gomd2svg` in all `.go` files.

---

### Task 4: Rename root package

**Files:**
- Modify: `mermaid.go`, `options.go`, `mermaid_test.go`, `mermaid_bench_test.go`, `mermaid_directive_test.go`, `mermaid_golden_test.go`

**Step 1:** Change `package mermaid` to `package gomd2svg` in all root-level .go files.

**Step 2:** Update the import alias in `cmd/gomd2svg/main.go` from `mermaid "..."` to direct import (package name matches).

**Step 3:** Update all references from `mermaid.Render` / `mermaid.RenderWithOptions` / `mermaid.RenderWithTiming` / `mermaid.Options` to `gomd2svg.Render` etc. in cmd/.

---

### Task 5: Rename root-level files

**Step 1:** Rename files

```bash
git mv mermaid.go gomd2svg.go
git mv mermaid_test.go gomd2svg_test.go
git mv mermaid_bench_test.go gomd2svg_bench_test.go
git mv mermaid_directive_test.go gomd2svg_directive_test.go
git mv mermaid_golden_test.go gomd2svg_golden_test.go
```

---

### Task 6: Rename CLI binary directory and update text

**Files:**
- Rename: `cmd/gomd2svg/` -> `cmd/gomd2svg/`
- Modify: `cmd/gomd2svg/main.go` (version string, usage text)
- Modify: `cmd/gomd2svg/main_test.go` (assertion string)

**Step 1:** Rename directory

```bash
git mv cmd/gomd2svg cmd/gomd2svg
```

**Step 2:** Update all `gomd2svg` strings in main.go and main_test.go to `gomd2svg`.

---

### Task 7: Update documentation

**Files:**
- Modify: all files in `docs/plans/` containing `gomd2svg`

**Step 1:** Bulk find-and-replace `gomd2svg` with `gomd2svg` in docs.

---

### Task 8: Verify build and tests

**Step 1:** Build

```bash
go build ./...
```

**Step 2:** Run all tests

```bash
go test ./...
```

**Step 3:** Run golden tests specifically

```bash
go test -run TestGolden -count=1
```

Expected: all 250 subtests pass.

---

### Task 9: Update Claude memory

**Files:**
- Modify: `/Users/jamesainslie/.claude/projects/-Volumes-Development-gomd2svg/memory/MEMORY.md`

**Step 1:** Replace all `gomd2svg` references with `gomd2svg`.

---

### Task 10: Commit and push

**Step 1:** Stage all changes, commit, push.

```bash
git add -A
git commit -m "rename: gomd2svg to gomd2svg"
git push
```
