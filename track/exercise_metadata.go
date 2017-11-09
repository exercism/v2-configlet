package track

// ExerciseMetadata contains metadata about an implemented exercise.
// It's listed in the config in the order that the exercise will be
// delivered by the API.
type ExerciseMetadata struct {
	UUID         string
	Slug         string
	Difficulty   int
	Topics       []string
	UnlockedBy   string `json:"unlocked_by"`
	IsCore       bool   `json:"core"`
	IsDeprecated bool   `json:"deprecated"`
}

// HasUUID checks if UUID is defined
func (e ExerciseMetadata) HasUUID() bool { return e.UUID != "" }

// ExerciseMetadataList is a slice of ExerciseMetadata to define functions
type ExerciseMetadataList []ExerciseMetadata

// Fold iterates over slice, runs given function and collects slugs
func (els ExerciseMetadataList) Fold(isValid func(ExerciseMetadata) bool) (valid []string, invalid []string) {
	for _, e := range els {
		if isValid(e) {
			valid = append(valid, e.Slug)
		} else {
			invalid = append(invalid, e.Slug)
		}
	}
	return
}
