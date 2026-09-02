package config

import "os"

// InstanceDecisionsGroup returns a consumer group unique to this instance.
// Every proxy replica must see EVERY decision (broadcast), so each instance
// gets its own group — a shared group would split decisions across
// replicas. In Kubernetes the hostname is the pod name.
func (c Config) InstanceDecisionsGroup() string {
	hostname, err := os.Hostname()
	if err != nil {
		return c.DecisionsGroup
	}
	return c.DecisionsGroup + "-" + hostname
}
