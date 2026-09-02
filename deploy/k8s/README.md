# Kubernetes manifests (Docker & K8s phase — priority 3)

One Helm chart per service; environment differences live in
`values-<env>.yaml` only. Every Deployment ships the full checklist from
docs/standards/infra.md: readiness + liveness probes on separate endpoints,
envFrom ConfigMap + Secret, resource requests and memory limits,
runAsNonRoot + readOnlyRootFilesystem + dropped capabilities, and
`terminationGracePeriodSeconds` aligned with the app's shutdown timeout.

Local target: kind cluster (installed). Charts arrive in the K8s phase after
the Compose stack works end to end.
