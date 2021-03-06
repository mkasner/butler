package butler

import (
	"os"
	"runtime"
	"strings"
)

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	// fmt.Println(os.Getenv("HOME"))
	return os.Getenv("HOME")
}

//Trims whitespace and new line chars
func strip(v string) string {
	return strings.TrimSpace(strings.Trim(v, "\n"))
}
