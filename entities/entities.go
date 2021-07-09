package entities

type User struct {
	ID         string   `datastore:"id"`
	Email      string   `datastore:"email"`
	Password   string   `datastore:"password"`
	ProjectIDs []string `datastore:"projectIds"`
}

type Session struct {
	ID        string   `datastore:"id"`
	ProjectID string   `datastore:"projectId"`
	SpecIDs   []string `datastore:"specIds"`
	Start     int64    `datastore:"start"`
	End       int64    `datastore:"end"`
}

type SessionWithSpecs struct {
	ID        string
	ProjectID string
	Specs     []Spec
	Start     int64
	End       int64
}

type Project struct {
	ID            string   `datastore:"id"`
	Name          string   `datastore:"name"`
	SessionIDs    []string `datastore:"sessionIds"`
	LatestSession string   `datastore:"latestSession"`
}

type ProjectFull struct {
	Sessions      []SessionWithSpecs
	LatestSession string
}

type Spec struct {
	ID                string `datastore:"id"`
	SessionID         string `datastore:"sessionId"`
	FilePath          string `datastore:"filePath"`
	Tests             []string
	EstimatedDuration int64  `datastore:"estimatedDuration"`
	Start             int64  `datastore:"start"`
	End               int64  `datastore:"end"`
	Passed            bool   `datastore:"passed"`
	AssignedTo        string `datastore:"assignedTo"`
}
