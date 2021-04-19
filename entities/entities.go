package entities

type User struct {
	ID         string
	Username   string
	Password   string
	ProjectIDs []string
}

type Session struct {
	ID        string
	ProjectID string
	Backlog   []Spec
	Start     int64
	End       int64
}

type Project struct {
	ID            string
	Name          string
	SessionIDs    []string
	LatestSession string
}

type ProjectFull struct {
	Sessions      []Session
	LatestSession string
}

type Spec struct {
	FilePath          string
	Tests             []string
	EstimatedDuration int64
	Start             int64
	End               int64
	AssignedTo        string
}
