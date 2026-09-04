package envx

import "fmt"

// Require fails when any named value is empty, reporting the first missing.
func Require(values map[string]string) error {
	for name, v := range values {
		if v == "" {
			return fmt.Errorf("config: %s is required", name)
		}
	}

	return nil
}
