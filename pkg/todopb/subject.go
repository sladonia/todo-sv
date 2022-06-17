package todopb

import "fmt"

const (
	ServiceName              = "todo-sv"
	ProjectDomainName        = "project"
	UserProjectEventsSubject = "user-project-events"
)

func NewProjectSubject(eventType, id string) string {
	return fmt.Sprintf("%s.%s.%s.%s", ServiceName, ProjectDomainName, eventType, id)
}

// example: todo-sv.user-project-events.PROJECT_CREATED.cao4dmp9d3pmus59pubg
func NewUserEventsSubject(eventType, userID string) string {
	return fmt.Sprintf("%s.%s.%s.%s", ServiceName, UserProjectEventsSubject, eventType, userID)
}
