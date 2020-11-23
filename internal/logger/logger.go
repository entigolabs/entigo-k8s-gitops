package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var Logger = log.New(os.Stderr, "", 0)

type Warning struct {
	Reason error
}

func (w *Warning) Error() string {
	warning := fmt.Sprintf("[warning] %s", w.Reason)
	return fmt.Sprintf("\x1b[36;1m%s\x1b[0m", warning)
}

type PrefixedError struct {
	Reason error
}

func (pe *PrefixedError) Error() string {
	error := fmt.Sprintf("[error] %s", pe.Reason)
	return fmt.Sprintf("\x1b[31;1m%s\x1b[0m", error)
}

func ChooseLogger(env string) *log.Logger {
	switch env {
	case "dev":
		return log.New(os.Stderr, "gitops: ", log.LstdFlags|log.Lshortfile)
	case "prod":
		return log.New(os.Stderr, "", 0)
	default:
		Logger.Println(&PrefixedError{errors.New("unsupported environment variable")})
		os.Exit(1)
		return nil
	}
}
