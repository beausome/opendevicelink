# Security policy

OpenDeviceLink is pre-alpha software and currently operates only on simulated devices. Do not use it to protect sensitive peripherals or data.

## Report a vulnerability

Until a public security address exists, do not open an issue containing exploit details. Contact the repository owner privately and include:

- the affected commit;
- a minimal reproduction;
- the impact and device classes involved;
- any suggested mitigation.

## Security boundaries

The current code tests pairing primitives, session transitions, exclusive leases, conservative device policy, and bounded framing. It does **not** yet provide:

- authenticated network discovery;
- encrypted device-traffic transport;
- operating-system driver isolation;
- protection from a compromised host;
- independent cryptographic review;
- physical-device support.

Any future real backend must keep HID, storage, and unknown classes disabled by default, release ownership on every failure path, and include adversarial protocol tests.

