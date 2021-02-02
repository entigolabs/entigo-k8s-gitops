package common

type Command int

const (
	RunCmd Command = iota
	UpdateCmd
	CopyCmd
	DeleteCmd
)

type LoggingLevel int

const (
	DevLoggingLvl LoggingLevel = iota
	ProdLoggingLvl
)
