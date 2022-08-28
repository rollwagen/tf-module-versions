package tf

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

type Module struct {
	Name     string
	Location struct {
		FileName string
		Line     int
	}
	UsedVersion      string
	AvailableVersion string
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
