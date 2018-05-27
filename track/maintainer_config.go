package track

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// MaintainerConfig contains the list of current and previous maintainers.
// The files is used both to manage the GitHub maintainer team, as well
// as to configure the display values for each maintainer on the Exercism
// website.
type MaintainerConfig struct {
	DocsURL     string       `json:"docs_url"`
	Maintainers []Maintainer `json:"maintainers"`
}

// Maintainer contains data about a track maintainer.
type Maintainer struct {
	Username      string  `json:"github_username"`
	Alumnus       bool    `json:"alumnus"`
	ShowOnWebsite bool    `json:"show_on_website"`
	Name          *string `json:"name"`
	LinkText      *string `json:"link_text"`
	LinkURL       *string `json:"link_url"`
	AvatarURL     *string `json:"avatar_url"`
	Bio           *string `json:"bio"`
}

// NewMaintainerConfig reads the maintainer config file, if present.
func NewMaintainerConfig(path string) (MaintainerConfig, error) {
	mc := MaintainerConfig{}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return mc, nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return mc, err
	}
	err = json.Unmarshal(bytes, &mc)
	if err != nil {
		return mc, fmt.Errorf("invalid config %s -- %s", path, err.Error())
	}
	return mc, nil
}
