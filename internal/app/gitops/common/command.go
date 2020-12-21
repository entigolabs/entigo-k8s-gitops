package common

type Command int

const (
	RunCmd Command = iota
	UpdateCmd
	CopyCmd
)
