package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseClassSimple(t *testing.T) {
	input := `classDiagram
class Animal {
  +String name
  +int age
  +makeSound() String
  +move(int distance) void
}`
	out, err := parseClass(input)
	if err != nil {
		t.Fatalf("parseClass() error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Class {
		t.Errorf("Kind = %v, want Class", g.Kind)
	}
	cm, ok := g.Members["Animal"]
	if !ok {
		t.Fatalf("Members map missing Animal")
	}
	if len(cm.Attributes) != 2 {
		t.Errorf("Attributes = %d, want 2", len(cm.Attributes))
	}
	if len(cm.Methods) != 2 {
		t.Errorf("Methods = %d, want 2", len(cm.Methods))
	}
	// Check method return type
	found := false
	for _, m := range cm.Methods {
		if m.Name == "makeSound" {
			found = true
			if m.Type != "String" {
				t.Errorf("makeSound return type = %q, want %q", m.Type, "String")
			}
		}
	}
	if !found {
		t.Error("method makeSound not found")
	}
}

func TestParseClassRelationships(t *testing.T) {
	tests := []struct {
		name  string
		arrow string
	}{
		{"inheritance", "Animal <|-- Duck"},
		{"composition", "Car *-- Engine"},
		{"aggregation", "Pond o-- Duck"},
		{"association", "Cat --> Owner"},
		{"dependency", "Cat ..> Food"},
		{"realization", "Animal ..|> Flying"},
		{"solid link", "Cat -- Dog"},
		{"dashed link", "Cat .. Dog"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := "classDiagram\n" + tt.arrow
			out, err := parseClass(input)
			if err != nil {
				t.Fatalf("parseClass() error: %v", err)
			}
			if len(out.Graph.Edges) != 1 {
				t.Errorf("Edges = %d, want 1", len(out.Graph.Edges))
			}
		})
	}
}

func TestParseClassAnnotation(t *testing.T) {
	input := `classDiagram
<<interface>> Shape`
	out, err := parseClass(input)
	if err != nil {
		t.Fatalf("parseClass() error: %v", err)
	}
	ann, ok := out.Graph.Annotations["Shape"]
	if !ok {
		t.Fatal("Annotations missing Shape")
	}
	if ann != "interface" {
		t.Errorf("Annotation = %q, want %q", ann, "interface")
	}
}

func TestParseClassDirection(t *testing.T) {
	input := `classDiagram
direction LR`
	out, err := parseClass(input)
	if err != nil {
		t.Fatalf("parseClass() error: %v", err)
	}
	if out.Graph.Direction != ir.LeftRight {
		t.Errorf("Direction = %v, want LeftRight", out.Graph.Direction)
	}
}

func TestParseClassVisibility(t *testing.T) {
	input := `classDiagram
class Foo {
  +publicAttr int
  -privateAttr int
  #protectedAttr int
  ~packageAttr int
}`
	out, err := parseClass(input)
	if err != nil {
		t.Fatalf("parseClass() error: %v", err)
	}
	cm, ok := out.Graph.Members["Foo"]
	if !ok {
		t.Fatalf("Members missing Foo")
	}
	if len(cm.Attributes) != 4 {
		t.Fatalf("Attributes = %d, want 4", len(cm.Attributes))
	}
	expected := []ir.Visibility{ir.VisPublic, ir.VisPrivate, ir.VisProtected, ir.VisPackage}
	for i, attr := range cm.Attributes {
		if attr.Visibility != expected[i] {
			t.Errorf("attr[%d] %q Visibility = %v, want %v", i, attr.Name, attr.Visibility, expected[i])
		}
	}
}

func TestParseClassCardinality(t *testing.T) {
	input := `classDiagram
Customer "1" --> "*" Order : places`
	out, err := parseClass(input)
	if err != nil {
		t.Fatalf("parseClass() error: %v", err)
	}
	if len(out.Graph.Edges) != 1 {
		t.Fatalf("Edges = %d, want 1", len(out.Graph.Edges))
	}
	edge := out.Graph.Edges[0]
	if edge.StartLabel == nil || *edge.StartLabel != "1" {
		t.Errorf("StartLabel = %v, want %q", edge.StartLabel, "1")
	}
	if edge.EndLabel == nil || *edge.EndLabel != "*" {
		t.Errorf("EndLabel = %v, want %q", edge.EndLabel, "*")
	}
	if edge.Label == nil || *edge.Label != "places" {
		t.Errorf("Label = %v, want %q", edge.Label, "places")
	}
}

func TestParseClassNamespace(t *testing.T) {
	input := `classDiagram
namespace BaseShapes {
  class Triangle
  class Rectangle
}`
	out, err := parseClass(input)
	if err != nil {
		t.Fatalf("parseClass() error: %v", err)
	}
	if len(out.Graph.Namespaces) != 1 {
		t.Fatalf("Namespaces = %d, want 1", len(out.Graph.Namespaces))
	}
	ns := out.Graph.Namespaces[0]
	if ns.Name != "BaseShapes" {
		t.Errorf("Namespace name = %q, want %q", ns.Name, "BaseShapes")
	}
	if len(ns.Classes) != 2 {
		t.Errorf("Namespace classes = %d, want 2", len(ns.Classes))
	}
}

func TestParseClassGeneric(t *testing.T) {
	input := `classDiagram
class List~T~`
	out, err := parseClass(input)
	if err != nil {
		t.Fatalf("parseClass() error: %v", err)
	}
	if _, ok := out.Graph.Nodes["List"]; !ok {
		t.Error("expected node List")
	}
}

func TestParseClassColon(t *testing.T) {
	input := `classDiagram
Animal : +int age`
	out, err := parseClass(input)
	if err != nil {
		t.Fatalf("parseClass() error: %v", err)
	}
	cm, ok := out.Graph.Members["Animal"]
	if !ok {
		t.Fatalf("Members missing Animal")
	}
	if len(cm.Attributes) != 1 {
		t.Fatalf("Attributes = %d, want 1", len(cm.Attributes))
	}
	if cm.Attributes[0].Name != "age" {
		t.Errorf("Attribute name = %q, want %q", cm.Attributes[0].Name, "age")
	}
}

func TestParseClassNote(t *testing.T) {
	input := `classDiagram
class Animal
note for Animal "This is a note"`
	out, err := parseClass(input)
	if err != nil {
		t.Fatalf("parseClass() error: %v", err)
	}
	if len(out.Graph.Notes) != 1 {
		t.Fatalf("Notes = %d, want 1", len(out.Graph.Notes))
	}
	if out.Graph.Notes[0].Text != "This is a note" {
		t.Errorf("Note text = %q, want %q", out.Graph.Notes[0].Text, "This is a note")
	}
	if out.Graph.Notes[0].Target != "Animal" {
		t.Errorf("Note target = %q, want %q", out.Graph.Notes[0].Target, "Animal")
	}
}

func TestParseClassClassifier(t *testing.T) {
	input := `classDiagram
class Foo {
  +abstractMethod()* void
  +staticMethod()$ void
}`
	out, err := parseClass(input)
	if err != nil {
		t.Fatalf("parseClass() error: %v", err)
	}
	cm, ok := out.Graph.Members["Foo"]
	if !ok {
		t.Fatalf("Members missing Foo")
	}
	if len(cm.Methods) != 2 {
		t.Fatalf("Methods = %d, want 2", len(cm.Methods))
	}
	foundAbstract := false
	foundStatic := false
	for _, m := range cm.Methods {
		if m.Name == "abstractMethod" && m.Classifier == ir.ClassifierAbstract {
			foundAbstract = true
		}
		if m.Name == "staticMethod" && m.Classifier == ir.ClassifierStatic {
			foundStatic = true
		}
	}
	if !foundAbstract {
		t.Error("expected abstractMethod with ClassifierAbstract")
	}
	if !foundStatic {
		t.Error("expected staticMethod with ClassifierStatic")
	}
}
