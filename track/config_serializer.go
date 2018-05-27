package track

type ConfigSerializer interface {
	Read(string) error
	ToJSON() ([]byte, error)
}
