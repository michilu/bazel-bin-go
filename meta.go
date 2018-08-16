package main

import (
	"v/cmd"
	"v/errs"
	"v/log"
	"v/meta"
)

// override by ldflags.
var (
	build  string
	hash   string
	serial string
)

func init() {
	const op = "main.init"
	err := meta.Set(&meta.Meta{
		Build:  build,
		Hash:   hash,
		Name:   name,
		SemVer: semVer,
		Serial: serial,
	})
	if err != nil {
		log.Logger().Fatal().
			Str("op", op).
			Err(&errs.Error{Op: op, Err: err}).
			Msg("error")
	}
	cmd.Init(ns)
}
