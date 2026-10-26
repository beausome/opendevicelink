package device

import "testing"

func TestDefaultPolicyBlocksDangerousClasses(t *testing.T) {
	policy := DefaultPolicy()
	for _, class := range []Class{ClassHID, ClassMassStorage, ClassUnknown} {
		decision := policy.Evaluate(Descriptor{Class: class})
		if decision.Allowed {
			t.Fatalf("class %q unexpectedly allowed", class)
		}
	}
}

func TestDefaultPolicyAllowsCamera(t *testing.T) {
	if decision := DefaultPolicy().Evaluate(Descriptor{Class: ClassCamera}); !decision.Allowed {
		t.Fatalf("camera unexpectedly blocked: %s", decision.Reason)
	}
}

func TestDescriptorValidation(t *testing.T) {
	valid := Descriptor{ID: "cam-1", Name: "Camera", VendorID: 1, ProductID: 2}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid descriptor rejected: %v", err)
	}
	invalid := valid
	invalid.ID = ""
	if err := invalid.Validate(); err == nil {
		t.Fatal("descriptor without ID accepted")
	}
}
