package mermaid

import (
	"os"
	"testing"
)

func BenchmarkRenderSimple(b *testing.B) {
	input := "flowchart LR; A-->B-->C"
	b.ReportAllocs()
	for b.Loop() {
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
	for b.Loop() {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderClassSimple(b *testing.B) {
	input := `classDiagram
    class Animal {
        +String name
        +int age
        +isMammal() bool
    }
    class Dog {
        +String breed
        +bark() void
    }
    Animal <|-- Dog`
	b.ReportAllocs()
	for b.Loop() {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderStateDiagram(b *testing.B) {
	input := `stateDiagram-v2
    [*] --> Still
    Still --> Moving
    Moving --> Crash
    Crash --> [*]`
	b.ReportAllocs()
	for b.Loop() {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderERDiagram(b *testing.B) {
	input := `erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ LINE-ITEM : contains
    CUSTOMER {
        string name
        int id PK
    }`
	b.ReportAllocs()
	for b.Loop() {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func readBenchFixture(b *testing.B, name string) string {
	b.Helper()
	data, err := os.ReadFile("testdata/fixtures/" + name)
	if err != nil {
		b.Fatalf("readBenchFixture(%q): %v", name, err)
	}
	return string(data)
}

func BenchmarkRenderSequenceSimple(b *testing.B) {
	input := readBenchFixture(b, "sequence-simple.mmd")
	b.ReportAllocs()
	for b.Loop() {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderSequenceComplex(b *testing.B) {
	input := readBenchFixture(b, "sequence-full.mmd")
	b.ReportAllocs()
	for b.Loop() {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderKanban(b *testing.B) {
	input := readBenchFixture(b, "kanban-basic.mmd")
	b.ReportAllocs()
	for b.Loop() {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderPacket(b *testing.B) {
	input := readBenchFixture(b, "packet-tcp.mmd")
	b.ReportAllocs()
	for b.Loop() {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
