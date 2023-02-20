package config

import (
	"os"
	"runtime"
)

func Setup() string {
	if runtime.GOOS == "windows" {
		dir := "C:\\tonic\\storage\\"

		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				panic(err)
			}
		} else if err != nil {
			panic(err)
		}
		return dir
	} else {
		// This is the linux path creation to determine storage locations
		dir := "/src/tonic/storage/"

		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				panic(err)
			}
		} else if err != nil {
			panic(err)
		}

		return dir
	}
}
