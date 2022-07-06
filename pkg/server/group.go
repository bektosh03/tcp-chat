package server

type Group struct {
	clients []*Client
}

func NewGroup(clients ...*Client) *Group {
	if clients == nil {
		clients = make([]*Client, 0)
	}
	return &Group{
		clients: clients,
	}
}
