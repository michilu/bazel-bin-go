package meta

import (
	"fmt"
	"runtime"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	m *meta
)

type (
	meta struct {
		Build   time.Time `yaml:",omitempty"`
		Hash    string    `yaml:",omitempty"`
		Name    string    `yaml:",omitempty"`
		SemVer  string    `yaml:",omitempty"`
		Serial  string    `yaml:",omitempty"`
		Runtime *runTime  `yaml:",omitempty"`
	}
	runTime struct {
		Version string `yaml:",omitempty"`
		Arch    string `yaml:",omitempty"`
		Os      string `yaml:",omitempty"`
	}
)

func (m meta) String() string {
	o, err := yaml.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return string(o)
}

func init() {
	m = &meta{
		Name:   name,
		Hash:   hash,
		SemVer: semVer,
		Serial: serial,
		Runtime: &runTime{
			Version: runtime.Version(),
			Arch:    runtime.GOARCH,
			Os:      runtime.GOOS,
		},
	}
	m.Build, _ = time.Parse(buildFmt, build)
}

// Get returns a fmt.Stringer.
func Get() fmt.Stringer {
	return *m
}

// Name returns a name.
func Name() string {
	return name
}
