package mermaid

import "testing"

func BenchmarkRenderSimple(b *testing.B) {
	input := "flowchart LR; A-->B-->C"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderMedium(b *testing.B) {
	input := `flowchart TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Process]
    B -->|No| D[Cancel]
    C --> E[End]
    D --> E
    E --> F[Cleanup]
    F --> G[Done]`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
