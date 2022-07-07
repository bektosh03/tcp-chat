package server

import (
	"fmt"
	"log"
)

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

func (g Group) chat(client *Client, msg string) {
	for _, c := range g.clients {
		if c.id == client.id {
			continue
		}
		if err := c.message(fmt.Sprintf("%s: %s", client.username, msg)); err != nil {
			log.Printf("failed to message client - %v: %v\n", c, err)
		}
	}
}
