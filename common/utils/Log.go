package utils

import (
	"fmt"
	"github.com/simpletonDL/GoGames/common/settings"
)

func Log(msg string, args ...any) {
	if settings.Debug {
		fmt.Printf(msg, args)
	}
}
