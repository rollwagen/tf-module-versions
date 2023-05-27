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

func NewModule(name string, usedVersion, availableVersion, ref, fileName string, line int) (*Module, error) {
	for _, v := range []string{usedVersion, availableVersion} {
		_, err := newVersion(v)
		if err != nil {
			return nil, err
		}
	}

	m := Module{
		Name:             name,
		UsedVersion:      usedVersion,
		AvailableVersion: availableVersion,
		GitReference:     ref,
	}
	m.Location.FileName = fileName
	m.Location.Line = line

	return &m, nil
}

func newVersion(versionAsString string) (*version.Version, error) {
	versionToCreate := "0"
	if versionAsString != "nil" {
		versionToCreate = versionAsString
	}
	newVersion, err := version.NewVersion(versionToCreate)
	if err != nil {
		return nil, fmt.Errorf("'%s' is not a valid version string: %w", versionAsString, err)
	}
	return newVersion, nil
}

func (m Module) HasNewerVersion() bool {
	ver := func(s string) *version.Version {
		v, _ := newVersion(s)
		return v
	}

	usedVersion := "0"
	if m.UsedVersion != "nil" {
		usedVersion = m.UsedVersion
	}

	return ver(usedVersion).LessThan(ver(m.AvailableVersion))
}

func (m Module) HasSameVersion() bool {
	return m.AvailableVersion == m.UsedVersion
}

func (m Module) NewerVersion() string {
	versionLevel := []string{"MAJOR", "MINOR", "PATCH"}

	const versionUndefined = "UNDEFINED"

	if !m.HasNewerVersion() {
		return versionUndefined
	}

	used, _ := newVersion(m.UsedVersion)

	available, _ := newVersion(m.AvailableVersion)

	for i := range used.Segments() {

		if i > 2 || (i > len(used.Segments())-1 && i > len(available.Segments())) {
			return versionUndefined
		}

		if available.Segments()[i] > used.Segments()[i] {
			return versionLevel[i]
		}
	}

	return versionUndefined
}
