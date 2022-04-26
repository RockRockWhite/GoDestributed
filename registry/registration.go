package registry

type Registration struct {
	ServiceName ServiceName
	ServiceUrl  string
}

type ServiceName string

const (
	LogSerice = ServiceName("LogService")
)
