package handle

import (
	"runtime"
)

func init() {
	if runtime.GOOS == "windows" {
		lineBreakLength = 2
	} else {
		lineBreakLength = 1
	}

}
