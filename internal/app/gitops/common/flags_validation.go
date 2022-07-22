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
	if isPathSet, _ := f.isFlagSet(f.App.Path, "app-path"); isPathSet {
		return false
	}
	return true
}

func (f *Flags) validatePath() error {
	// TODO app-path validation
	return nil
}

func (f *Flags) validateTokenizedPath() error {
	if _, err := f.isFlagSet(f.App.Prefix, "app-prefix"); err != nil {
		return err
	}
	if _, err := f.isFlagSet(f.App.Namespace, "app-namespace"); err != nil {
		return err
	}
	if _, err := f.isFlagSet(f.App.Name, "app-name"); err != nil {
		return err
	}
	return nil
}

func (f *Flags) isFlagSet(flagValue string, flagType string) (bool, error) {
	if flagValue != "" {
		return true, nil
	}
	return false, errors.New(flagNotSetMsg(flagType))
}

func (f *Flags) validateUpdate() error {
	isNotificationFlagSet := f.Notification.URL != "" || f.Notification.RegistryUri != "" ||
		f.Notification.AuthToken != "" || f.Notification.Environment != ""
	if !isNotificationFlagSet {
		return nil
	}
	return f.validateNotificationFlags()
}

func (f *Flags) validateNotificationFlags() error {
	if f.Notification.URL == "" {
		return errors.New(notifyFlagNotSetMsg("notify-api-url"))
	}
	if f.Notification.RegistryUri == "" {
		return errors.New(notifyFlagNotSetMsg("notify-registry-uri"))
	}
	if f.Notification.AuthToken == "" {
		return errors.New(notifyFlagNotSetMsg("notify-auth-token"))
	}
	if f.Notification.Environment == "" {
		return errors.New(notifyFlagNotSetMsg("notify-env"))
	}
	return nil
}

func flagNotSetMsg(flagType string) string {
	return fmt.Sprintf("required flag '%s' not set", flagType)
}

func notifyFlagNotSetMsg(flagType string) string {
	return fmt.Sprintf("required notification flag '%s' not set", flagType)
}
