package ir

import "testing"

func TestVisibilitySymbol(t *testing.T) {
	tests := []struct {
		vis  Visibility
		want string
	}{
		{VisNone, ""},
		{VisPublic, "+"},
		{VisPrivate, "-"},
		{VisProtected, "#"},
		{VisPackage, "~"},
	}
	for _, tt := range tests {
		got := tt.vis.Symbol()
		if got != tt.want {
			t.Errorf("Visibility(%d).Symbol() = %q, want %q", tt.vis, got, tt.want)
		}
	}
}

func TestClassMembersGrouping(t *testing.T) {
	members := &ClassMembers{
		Attributes: []ClassMember{
			{Name: "id", Type: "int", Visibility: VisPrivate},
			{Name: "name", Type: "string", Visibility: VisPrivate},
		},
		Methods: []ClassMember{
			{Name: "getId", Type: "int", Params: "", IsMethod: true, Visibility: VisPublic},
			{Name: "setName", Type: "void", Params: "String name", IsMethod: true, Visibility: VisPublic},
		},
	}

	if len(members.Attributes) != 2 {
		t.Errorf("Attributes count = %d, want 2", len(members.Attributes))
	}
	if len(members.Methods) != 2 {
		t.Errorf("Methods count = %d, want 2", len(members.Methods))
	}

	// Verify attributes are not methods
	for _, a := range members.Attributes {
		if a.IsMethod {
			t.Errorf("attribute %q should not be a method", a.Name)
		}
	}

	// Verify methods are methods
	for _, m := range members.Methods {
		if !m.IsMethod {
			t.Errorf("method %q should be a method", m.Name)
		}
	}
}

func TestGraphClassFields(t *testing.T) {
	g := NewGraph()

	if g.Members == nil {
		t.Error("Members map should be initialized")
	}
	if g.Annotations == nil {
		t.Error("Annotations map should be initialized")
	}
	if g.Namespaces != nil {
		t.Error("Namespaces should be nil (zero-value slice)")
	}
	if g.Notes != nil {
		t.Error("Notes should be nil (zero-value slice)")
	}
}
