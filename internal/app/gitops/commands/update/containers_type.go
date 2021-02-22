package update

type containersType int

const (
	containers containersType = iota
	initContainers
)

func (ct containersType) string() string {
	return [...]string{"containers", "initContainers"}[ct]
}
