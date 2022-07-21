package common

type Command int

const (
	RunCmd Command = iota
	UpdateCmd
	CopyCmd
	DeleteCmd
	DeploymentNotificationCmd
	ArgoCDDeleteCmd
	ArgoCDGetCmd
	ArgoCDSyncCmd
	ArgoCDUpdateCmd
	VersionCmd
)

type LoggingLevel int

const (
	DevLoggingLvl LoggingLevel = iota
	ProdLoggingLvl
)

type DeploymentStrategy int

const (
	RollingUpdateStrategy DeploymentStrategy = iota
	RecreateStrategy
	UnspecifiedStrategy
)

func (d DeploymentStrategy) String() string {
	return [...]string{"RollingUpdate", "Recreate"}[d]
}
