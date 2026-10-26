package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/beausome/opendevicelink/internal/lease"
	"github.com/beausome/opendevicelink/internal/securepair"
	"github.com/beausome/opendevicelink/internal/session"
	"github.com/beausome/opendevicelink/internal/simulator"
)

const usage = `OpenDeviceLink — safely hand off USB devices between computers

Usage:
  opendevicelink demo    run a simulated encrypted handoff
  opendevicelink list    list simulated devices and safety decisions
  opendevicelink doctor  report host and backend capabilities
  opendevicelink help    show this help

Current status: simulation only. No real USB device is attached or redirected.`

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, out, errOut io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(out, usage)
		return nil
	}

	switch args[0] {
	case "demo":
		return runDemo(out)
	case "list":
		return runList(out)
	case "doctor":
		return runDoctor(out)
	case "help", "-h", "--help":
		fmt.Fprintln(out, usage)
		return nil
	default:
		fmt.Fprintln(errOut, usage)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runDemo(out io.Writer) error {
	backend := simulator.New()
	devices, err := backend.List()
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		return errors.New("simulator returned no devices")
	}

	selected := devices[0]
	if decision := simulator.DefaultPolicy().Evaluate(selected); !decision.Allowed {
		return fmt.Errorf("selected device is blocked: %s", decision.Reason)
	}

	code, err := securepair.GenerateCode()
	if err != nil {
		return fmt.Errorf("generate pairing code: %w", err)
	}
	host, err := securepair.NewPeer()
	if err != nil {
		return fmt.Errorf("create host identity: %w", err)
	}
	client, err := securepair.NewPeer()
	if err != nil {
		return fmt.Errorf("create client identity: %w", err)
	}
	hostKey, err := host.Derive(client.PublicKey(), selected.ID)
	if err != nil {
		return fmt.Errorf("derive host key: %w", err)
	}
	clientKey, err := client.Derive(host.PublicKey(), selected.ID)
	if err != nil {
		return fmt.Errorf("derive client key: %w", err)
	}
	if hostKey != clientKey {
		return errors.New("pairing produced mismatched session keys")
	}

	flow := session.New(selected.ID)
	if err := flow.BeginPairing(code); err != nil {
		return err
	}
	if err := flow.Authorize(code); err != nil {
		return err
	}

	leases := lease.NewManager(time.Now)
	deviceLease, err := leases.Acquire(selected.ID, "demo-client", 30*time.Second)
	if err != nil {
		return err
	}
	if err := backend.Attach(selected.ID); err != nil {
		return err
	}
	if err := flow.Attach(); err != nil {
		return err
	}

	fmt.Fprintln(out, "OpenDeviceLink simulated handoff")
	fmt.Fprintf(out, "1. discovered  %s (%s)\n", selected.Name, selected.ID)
	fmt.Fprintf(out, "2. paired      code %s; verification %s\n", code, securepair.VerificationCode(hostKey))
	fmt.Fprintln(out, "3. authorized  encrypted session keys match")
	fmt.Fprintf(out, "4. attached    exclusive lease %s\n", deviceLease.Token[:8])

	if err := backend.Detach(selected.ID); err != nil {
		return err
	}
	if err := flow.Detach(); err != nil {
		return err
	}
	if err := leases.Release(deviceLease.Token); err != nil {
		return err
	}
	fmt.Fprintln(out, "5. detached    lease released; device returned to host")
	fmt.Fprintln(out, "Result: PASS (simulation only; no physical USB traffic was sent)")
	return nil
}

func runList(out io.Writer) error {
	backend := simulator.New()
	devices, err := backend.List()
	if err != nil {
		return err
	}
	policy := simulator.DefaultPolicy()
	for _, item := range devices {
		decision := policy.Evaluate(item)
		status := "allowed"
		if !decision.Allowed {
			status = "blocked"
		}
		fmt.Fprintf(out, "%s  %04x:%04x  %-12s  %-7s  %s\n", item.ID, item.VendorID, item.ProductID, item.Class, status, decision.Reason)
	}
	return nil
}

func runDoctor(out io.Writer) error {
	fmt.Fprintln(out, "OpenDeviceLink capability report")
	fmt.Fprintf(out, "Host:             %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintln(out, "Discovery:        simulated")
	fmt.Fprintln(out, "Pairing crypto:   X25519 + HMAC-SHA-256 key derivation")
	fmt.Fprintln(out, "Device transport: simulated")
	fmt.Fprintln(out, "Real USB attach:  unavailable in this foundation release")
	fmt.Fprintln(out, "Status:           ready for protocol development, not hardware use")
	return nil
}
