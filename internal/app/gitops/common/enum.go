package common

type Command int

const (
	RunCmd Command = iota
	UpdateCmd
	CopyCmd
)

type LoggingLevel int

const (
	DevLoggingLvl LoggingLevel = iota
	ProdLoggingLvl
)
