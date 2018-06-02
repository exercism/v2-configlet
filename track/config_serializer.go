package track

type ConfigSerializer interface {
	NewConfigFromFile(string) error
	ToJSON() ([]byte, error)
}
