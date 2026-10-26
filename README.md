# OpenDeviceLink

OpenDeviceLink is an experimental foundation for handing a USB device from one computer to another with a short pairing code, encrypted session setup, conservative safety defaults, and exclusive device leases.

> **Current status: simulation only.** This repository does not yet discover, detach, attach, or transport traffic for physical USB devices. The demo exercises the protocol-independent state, security, framing, and lease logic so future operating-system backends can be added without pretending they already work.

## Why this exists

Remote USB software is often tied to one operating system, depends on manual networking, or exposes low-level tools directly to users. OpenDeviceLink is testing a narrower product promise: discover a nearby device, verify a six-digit code on both computers, grant one receiver a time-limited lease, and return the device safely when the session ends.

## What works today

- deterministic simulated USB discovery;
- conservative policy that blocks keyboards and storage by default;
- X25519 shared-secret agreement and context-bound HMAC-SHA-256 key derivation;
- six-digit human verification code derived from the session key;
- enforced `discovered → pairing → authorized → attached → detached` lifecycle;
- exclusive, expiring leases;
- bounded binary frames that reject malformed or oversized input;
- Windows, Linux, and macOS CI for the platform-independent core.

## Quick start

Go 1.27 or newer is required.

```sh
go test ./...
go run ./cmd/opendevicelink doctor
go run ./cmd/opendevicelink list
go run ./cmd/opendevicelink demo
```

The demo ends with this explicit boundary:

```text
Result: PASS (simulation only; no physical USB traffic was sent)
```

## Architecture

```text
CLI
 ├─ simulator backend ── device descriptors and attach lifecycle
 ├─ safety policy ────── default-deny for HID, storage, and unknown classes
 ├─ secure pairing ───── X25519 + context-bound HMAC-SHA-256
 ├─ session state ────── legal transitions only
 ├─ lease manager ────── one active receiver per device
 └─ bounded protocol ─── typed, length-prefixed frames
```

The operating-system boundary is intentionally absent from the first release. Real backends must satisfy the same policy, session, lease, and framing tests before they are described as supported.

## Safety model

- A device is never attachable before pairing and authorization.
- Only one unexpired lease can own a device.
- HID and mass-storage classes are blocked by default because remote input injection and filesystem corruption have unusually high consequences.
- Unknown classes fail closed.
- Frames are size-limited before payload allocation.
- Session keys are derived from ephemeral X25519 keys and a device-specific context.

This is not yet a security audit or a production guarantee. See [SECURITY.md](SECURITY.md).

## Honest support matrix

| Capability | Windows | Linux | macOS |
| --- | --- | --- | --- |
| Core protocol tests | CI | CI | CI |
| Simulator demo | Yes | Yes | Yes |
| Physical discovery | Not implemented | Not implemented | Not implemented |
| Real USB attach | Not implemented | Not implemented | Not implemented |
| Traffic transport | Not implemented | Not implemented | Not implemented |

## Roadmap

1. Specify a transport-neutral USB request/response model and fuzz its decoder.
2. Prototype Linux discovery and ownership using documented kernel interfaces.
3. Prototype a Windows backend and measure driver-signing constraints.
4. Run hardware-in-the-loop tests for cameras and printers before enabling either class.
5. Add authenticated local-network discovery and an encrypted transport.
6. Investigate macOS only after documenting its entitlement and reduced-security limitations.

Hardware support will be tracked by device model, host OS, receiver OS, test evidence, and known failure modes—never by an unqualified checkmark.

## Contributing

Read [CONTRIBUTING.md](CONTRIBUTING.md). Small, test-backed changes are preferred. Please do not describe hardware as supported without reproducible hardware-in-the-loop evidence.

## License

MIT.


