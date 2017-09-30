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
