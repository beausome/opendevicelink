// Package device defines USB device metadata and conservative safety policy.
package device

import "fmt"

type Class string

const (
	ClassAudio       Class = "audio"
	ClassCamera      Class = "camera"
	ClassHID         Class = "hid"
	ClassMassStorage Class = "storage"
	ClassPrinter     Class = "printer"
	ClassUnknown     Class = "unknown"
)

type Descriptor struct {
	ID        string
	Name      string
	VendorID  uint16
	ProductID uint16
	Class     Class
}

func (d Descriptor) Validate() error {
	if d.ID == "" {
		return fmt.Errorf("device ID is required")
	}
	if d.Name == "" {
		return fmt.Errorf("device name is required")
	}
	if d.VendorID == 0 || d.ProductID == 0 {
		return fmt.Errorf("vendor and product IDs must be non-zero")
	}
	return nil
}

type Decision struct {
	Allowed bool
	Reason  string
}

type Policy struct {
	allowDangerous bool
}

func DefaultPolicy() Policy {
	return Policy{}
}

func PermissivePolicy() Policy {
	return Policy{allowDangerous: true}
}

func (p Policy) Evaluate(d Descriptor) Decision {
	switch d.Class {
	case ClassHID:
		if !p.allowDangerous {
			return Decision{Reason: "human-interface devices are blocked by default to reduce input-injection risk"}
		}
	case ClassMassStorage:
		if !p.allowDangerous {
			return Decision{Reason: "mass-storage devices are blocked by default to reduce data-loss risk"}
		}
	case ClassUnknown:
		return Decision{Reason: "unknown device classes require an explicit policy"}
	}
	return Decision{Allowed: true, Reason: "class permitted by current policy"}
}
