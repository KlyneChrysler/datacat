package config

import "os"

// InstanceDecisionsGroup gives each replica its own group so all see every decision.
func (c Config) InstanceDecisionsGroup() string {
	hostname, err := os.Hostname()
	if err != nil {
		return c.DecisionsGroup
	}

	return c.DecisionsGroup + "-" + hostname
}
