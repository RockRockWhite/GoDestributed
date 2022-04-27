package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceUrl       string
	RequiredServices []ServiceName
	ServiceUpdateUrl string
}

type ServiceName string

const (
	LogSerice     = ServiceName("LogService")
	GradesService = ServiceName("GradesService")
)

type patchEntry struct {
	ServiceName ServiceName
	Url         string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
