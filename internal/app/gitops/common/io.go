package common

import (
	"os"
)

func GetFileInput(fileName string) []byte {
	input, err := os.ReadFile(fileName)
	if err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
	return input
}

func OverwriteFile(fileName string, data []byte) {
	if err := os.WriteFile(fileName, data, 0644); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
}
