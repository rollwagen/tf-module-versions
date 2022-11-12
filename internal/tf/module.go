package tf

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-version"
)

type Module struct {
	Name     string `json:"name"`
	Location struct {
		FileName string `json:"fileName"`
		Line     int    `json:"line"`
	} `json:"location"`
	UsedVersion      string `json:"usedVersion"`
	AvailableVersion string `json:"availableVersion"`
	GitReference     string `json:"gitRef"`
}

func (m Module) MarshalJSON() ([]byte, error) {
	type baseModule Module
	// Marshal baseModule with an extension
	return json.Marshal(
		struct {
			baseModule
			HasNewerVersion bool `json:"hasNewerVersion"`
		}{
			baseModule(m),
			m.HasNewerVersion(),
		},
	)
}

func NewModule(name string, usedVersion, availableVersion, fileName string, line int) (*Module, error) {
	for _, v := range []string{usedVersion, availableVersion} {
		_, err := version.NewVersion(v)
		if err != nil {
			return nil, fmt.Errorf("'%s' is not a valid version string: %w", v, err)
		}
	}

	m := Module{
		Name:             name,
		UsedVersion:      usedVersion,
		AvailableVersion: availableVersion,
	}
	m.Location.FileName = fileName
	m.Location.Line = line

	return &m, nil
}

func (m Module) HasNewerVersion() bool {
	ver := func(s string) *version.Version {
		v, _ := version.NewVersion(s)
		return v
	}
	return ver(m.UsedVersion).LessThan(ver(m.AvailableVersion))
}

func (m Module) HasSameVersion() bool {
	return m.AvailableVersion == m.UsedVersion
}
