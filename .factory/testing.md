test-gate: go test ./...

<!-- halo:test-layers -->
{"layers": [
  {"id": "unit", "class": "surrogate", "command": "go test ./...", "required": true, "runBy": "canonical-gate"}
],
"productionPathNotApplicable": "The repository ships a CLI binary (cmd/calc) but no automated test exercises the compiled binary end-to-end; correctness of the binary was verified manually against all acceptance criteria during delivery.",
"policyNotApplicable": "No licence, compliance, or dependency-policy checks exist for this repository; it uses only the Go standard library."
}
