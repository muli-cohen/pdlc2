# Testing

test-gate: go test ./...

<!-- halo:test-layers -->
```json
{
  "layers": [
    {"id": "unit", "class": "surrogate", "command": "go test ./...", "required": true, "runBy": "canonical-gate"},
    {"id": "cli", "class": "production-path", "command": "go test ./cmd/joke/...", "required": true, "runBy": "canonical-gate"}
  ],
  "policyNotApplicable": "Standard library only: no third-party dependencies, licence, or compliance checks apply."
}
```
