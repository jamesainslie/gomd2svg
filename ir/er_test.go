package ir

import "testing"

func TestAttributeKeyString(t *testing.T) {
	tests := []struct {
		key  AttributeKey
		want string
	}{
		{KeyNone, ""},
		{KeyPrimary, "PK"},
		{KeyForeign, "FK"},
		{KeyUnique, "UK"},
	}
	for _, tt := range tests {
		got := tt.key.String()
		if got != tt.want {
			t.Errorf("AttributeKey(%d).String() = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestEntityDisplayName_WithLabel(t *testing.T) {
	e := &Entity{
		ID:    "customers",
		Label: "Customers",
	}
	got := e.DisplayName()
	if got != "Customers" {
		t.Errorf("DisplayName() = %q, want %q", got, "Customers")
	}
}

func TestEntityDisplayName_WithoutLabel(t *testing.T) {
	e := &Entity{
		ID:    "customers",
		Label: "",
	}
	got := e.DisplayName()
	if got != "customers" {
		t.Errorf("DisplayName() = %q, want %q", got, "customers")
	}
}
