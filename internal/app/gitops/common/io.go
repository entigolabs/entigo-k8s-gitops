package common

import "io/ioutil"

func GetFileInput(fileName string) []byte {
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
	return input
}

func OverwriteFile(fileName string, output string) {
	if err := ioutil.WriteFile(fileName, []byte(output), 0644); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
}
