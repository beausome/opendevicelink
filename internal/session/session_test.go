package session

import "testing"

func TestHappyPath(t *testing.T) {
	flow := New("camera-1")
	steps := []struct {
		name string
		do   func() error
		want State
	}{
		{"pair", func() error { return flow.BeginPairing("123456") }, Pairing},
		{"authorize", func() error { return flow.Authorize("123456") }, Authorized},
		{"attach", flow.Attach, Attached},
		{"detach", flow.Detach, Detached},
	}
	for _, step := range steps {
		if err := step.do(); err != nil {
			t.Fatalf("%s: %v", step.name, err)
		}
		if flow.State() != step.want {
			t.Fatalf("%s state = %s, want %s", step.name, flow.State(), step.want)
		}
	}
}

func TestInvalidTransitionsAreRejected(t *testing.T) {
	flow := New("camera-1")
	if err := flow.Attach(); err == nil {
		t.Fatal("attach before authorization accepted")
	}
	if err := flow.BeginPairing("123456"); err != nil {
		t.Fatalf("begin pairing: %v", err)
	}
	if err := flow.Authorize("654321"); err == nil {
		t.Fatal("incorrect pairing code accepted")
	}
	if flow.State() != Pairing {
		t.Fatalf("failed authorization changed state to %s", flow.State())
	}
}
