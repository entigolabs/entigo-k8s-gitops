package common

import (
	"errors"
	"fmt"
)

func (f *Flags) validate(cmd Command) error {
	if cmd == UpdateCmd {
		return f.validatePaths()
	}
	return nil
}

func (f *Flags) validatePaths() error {
	if !f.isTokenizedPath() {
		if err := f.validatePath(); err != nil {
			return err
		}
	} else {
		if err := f.validateTokenizedPath(); err != nil {
			return err
		}
	}
	return nil
}

func (f *Flags) isTokenizedPath() bool {
	if isPathSet, _ := f.isFlagSet(f.App.Path); isPathSet {
		return false
	}
	return true
}

func (f *Flags) validatePath() error {
	// TODO app-path validation
	return nil
}

func (f *Flags) validateTokenizedPath() error {
	if _, err := f.isFlagSet(f.App.Prefix); err != nil {
		return err
	}
	if _, err := f.isFlagSet(f.App.Namespace); err != nil {
		return err
	}
	if _, err := f.isFlagSet(f.App.Name); err != nil {
		return err
	}
	return nil
}

func (f *Flags) isFlagSet(flagValue string) (bool, error) {
	if flagValue != "" {
		return true, nil
	}
	return false, errors.New(fmt.Sprintf(`required flag "%s" not set`, f.App.Prefix))
}
