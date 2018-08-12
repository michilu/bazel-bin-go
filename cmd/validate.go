package cmd

import (
	"runtime"

	valid "github.com/asaskevich/govalidator"
)

func init() {
	valid.TagMap["filepath"] = isFilePath
}

func isFilePath(s string) bool {
	ok, os := valid.IsFilePath(s)
	if !ok {
		return false
	}
	switch os {
	case valid.Win:
		if runtime.GOOS != "windows" {
			return false
		}
	}
	return true
}
