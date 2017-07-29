package track

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type MaintainerConfig struct {
	Maintainers []Maintainer `json:"maintainers"`
	DocsURL     string       `json:"docs_url"`
}

type Maintainer struct {
	Username      string `json:"github_username"`
	ShowOnWebsite bool   `json:"show_on_website"`
	Alumnus       bool   `json:"alumnus"`
	Name          string `json:"name"`
	Bio           string `json:"bio"`
	LinkText      string `json:"link_text"`
	LinkURL       string `json:"link_url"`
	AvatarURL     string `json:"avatar_url"`
}

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
