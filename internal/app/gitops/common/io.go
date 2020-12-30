package common

import "io/ioutil"

func GetFileInput(fileName string) []byte {
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
	return input
}

func OverwriteFile(fileName string, data []byte) {
	if err := ioutil.WriteFile(fileName, data, 0644); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
}
