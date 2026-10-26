# Contributing

OpenDeviceLink values evidence over broad compatibility claims.

1. Create a focused branch.
2. Add or update tests with the change.
3. Run `go test ./...`, `go vet ./...`, and `go test -race ./...` where supported.
4. Keep the simulator and real-backend boundaries explicit.
5. Document security consequences and failure behavior.

A pull request that claims physical-device support should include the tested device IDs, host and receiver operating-system versions, repeatable commands, logs with secrets removed, and at least one forced-disconnect recovery test.

