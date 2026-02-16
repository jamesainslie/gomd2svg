# Rename gomd2svg to gomd2svg

## Goal

Rename the project from `gomd2svg` to `gomd2svg` across all code, docs, and infrastructure.

## Scope

| Component | Before | After |
|-----------|--------|-------|
| GitHub repo | `jamesainslie/gomd2svg` | `jamesainslie/gomd2svg` |
| Go module | `github.com/jamesainslie/gomd2svg` | `github.com/jamesainslie/gomd2svg` |
| Root package | `package mermaid` | `package gomd2svg` |
| CLI binary | `cmd/gomd2svg/` | `cmd/gomd2svg/` |
| CLI command | `gomd2svg render ...` | `gomd2svg render ...` |
| Public API | `mermaid.Render()` | `gomd2svg.Render()` |

## Unchanged

- Subpackage names: `ir`, `parser`, `layout`, `render`, `theme`, `config`, `textmetrics`
- Test fixtures (`testdata/fixtures/*.mmd`)
- Golden SVG files (`testdata/golden/*.svg`)
- Internal logic and architecture

## Execution Order

1. Rename GitHub repo (auto-redirects old URLs)
2. Update `go.mod` module path
3. Find-and-replace all import paths
4. Rename root package and update callers
5. Rename `cmd/gomd2svg/` directory and update CLI text
6. Update documentation
7. Verify build and tests
8. Update Claude memory
9. Commit and push

## Decisions

- **Package name**: `package gomd2svg` (not keeping `package mermaid`)
- **GitHub redirect**: GitHub auto-redirects old repo URLs, so existing links keep working
- **No version bump**: This is a rename, not a new major version
