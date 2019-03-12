package osbapi

type Broker struct {
	Name      string
	CreatedAt string
	URL       string
	Username  string
	Password  string
}

type Service struct {
	ID          string
	Name        string
	Description string
	BrokerName  string
}

type Plan struct {
	ID        string
	Name      string
	ServiceID string
}

type Instance struct {
	Name       string
	ServiceID  string
	PlanID     string
	BrokerName string
}
