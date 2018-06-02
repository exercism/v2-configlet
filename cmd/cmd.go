package cmd

// ConfigSerializer reads and serializes JSON configs.
type ConfigSerializer interface {
	NewConfigFromFile(string) error
	ToJSON() ([]byte, error)
}
