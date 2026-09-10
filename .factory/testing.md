# Testing

test-gate: go test ./...

<!-- halo:test-layers -->
```json
{
  "layers": [
    {"id": "unit", "class": "surrogate", "command": "go test ./...", "required": true, "runBy": "canonical-gate"}
  ],
  "productionPathNotApplicable": "This repository ships a small Go CLI whose behaviour is covered by package unit tests; there is no deployed artefact or end-to-end harness to exercise separately.",
  "policyNotApplicable": "Standard library only: no third-party dependencies, licence, or compliance checks apply."
}
```
