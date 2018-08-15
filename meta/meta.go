package meta

import (
	"fmt"
	"runtime"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/michilu/bazel-bin-go/errs"
	"github.com/michilu/bazel-bin-go/log"
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
	const op = "meta.init"
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

	if build == "" {
		return
	}
	t, err := time.Parse(buildFmt, build)
	if err != nil {
		log.Logger().Fatal().
			Str("op", op).
			Err(&errs.Error{Op: op, Err: err}).
			Msg("error")
	}
	m.Build = t

}

// Get returns a fmt.Stringer.
func Get() fmt.Stringer {
	return *m
}

// Name returns a name.
func Name() string {
	return name
}
