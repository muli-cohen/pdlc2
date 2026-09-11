# Testing

test-gate: go test ./...

<!-- halo:test-layers -->
```json
{
  "layers": [
    {"id": "unit", "class": "surrogate", "command": "go test ./...", "required": true, "runBy": "canonical-gate"},
    {"id": "cli", "class": "production-path", "command": "go test ./cmd/...", "required": true, "runBy": "canonical-gate"},
    {"id": "deps-policy", "class": "policy", "command": "[ $(grep -c '^require' go.mod) -eq 0 ]", "required": true, "runBy": "local"}
  ],
  "policyNotApplicable": "Licence and compliance checks are not applicable: the implementation uses Go standard library only, verified by the deps-policy layer."
}
```
