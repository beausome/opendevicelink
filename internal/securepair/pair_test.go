package securepair

import "testing"

func TestPeersDeriveMatchingKeys(t *testing.T) {
	alice, err := NewPeer()
	if err != nil {
		t.Fatalf("new alice: %v", err)
	}
	bob, err := NewPeer()
	if err != nil {
		t.Fatalf("new bob: %v", err)
	}
	aliceKey, err := alice.Derive(bob.PublicKey(), "device-1")
	if err != nil {
		t.Fatalf("alice derive: %v", err)
	}
	bobKey, err := bob.Derive(alice.PublicKey(), "device-1")
	if err != nil {
		t.Fatalf("bob derive: %v", err)
	}
	if aliceKey != bobKey {
		t.Fatal("derived keys differ")
	}
}

func TestContextSeparatesKeys(t *testing.T) {
	alice, _ := NewPeer()
	bob, _ := NewPeer()
	first, _ := alice.Derive(bob.PublicKey(), "device-1")
	second, _ := alice.Derive(bob.PublicKey(), "device-2")
	if first == second {
		t.Fatal("different contexts produced the same key")
	}
}

func TestGeneratedCodeIsValid(t *testing.T) {
	code, err := GenerateCode()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if !ValidCode(code) {
		t.Fatalf("invalid generated code %q", code)
	}
}
