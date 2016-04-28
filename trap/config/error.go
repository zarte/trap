package config

import (
	"github.com/raincious/trap/trap/core/types"
)

var (
	ErrParseInvalidItem *types.Error = types.NewError("Invalid value \"%s\" in option `%s`.")

	ErrInvalidCmdItem *types.Error = types.NewError("The number '%d' option in `Command` set '%s' is invalid.")
)
