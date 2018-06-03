package cmd

// ConfigSerializer reads and serializes JSON configs.
type ConfigSerializer interface {
	LoadFromFile(string) error
	ToJSON() ([]byte, error)
}
