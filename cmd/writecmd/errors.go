package writecmd

import (
	"errors"
)

var (
	SkipFileErr = errors.New("SkipFileErr")
	NoFileErr   = errors.New("NoFileErr")
)
