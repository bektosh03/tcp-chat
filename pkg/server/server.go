package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	groups map[string]*Group
}

func New() Server {
	return Server{
		groups: make(map[string]*Group),
	}
}

func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Println("Waiting for connections")
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("error while accepting conn: %v\n", err)
			continue
		}

		log.Println("Accepted new connection")

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	c, err := NewClient(conn)
	if err != nil {
		log.Printf("failed to initialize new client: %v\n", err)
		return
	}
	log.Printf("%s joined our server!\n", c.username)

	if err := c.message(fmt.Sprintf("Welcome %s!", c.username)); err != nil {
		log.Printf("error while messaging client: %v, %v", err, c)
	}

	for {
		in, err := c.readInput()
		if err != nil {
			log.Printf("error while reading client input: %v\n", err)
			continue
		}

		ins := strings.Fields(in)
		if len(ins) < 1 {
			log.Printf("empty input")
			continue
		}
		switch ins[0] {
		case "/groups":
			s.listGroups(c)
		case "/create":
			if len(ins) != 2 {
				c.message("ERR: /create command requires one argument - group's name")
				continue
			}
			s.createGroup(&c, ins[1])
		case "/join":
			if len(ins) != 2 {
				c.message("ERR: /join command requires one argument - group's name")
				continue
			}
			s.joinGroup(&c, ins[1])
		case "/exit":
			log.Printf("%s left server\n", c.username)
			return
		default:
			c.message("ERR: unrecognized command!")
		}
	}
}

func (s *Server) listGroups(c Client) {
	msg := strings.Builder{}
	msg.WriteString("Availabe groups:\n")

	for k := range s.groups {
		msg.WriteString(fmt.Sprintf("%s\n", k))
	}

	if err := c.message(msg.String()); err != nil {
		log.Printf("error while messaging client: %v - %v", err, c)
	}
}

func (s *Server) createGroup(c *Client, groupName string) {
	if _, ok := s.groups[groupName]; ok {
		c.message("ERR: this group already exists!")
		return
	}

	s.groups[groupName] = NewGroup(c)
	c.group = s.groups[groupName]
	c.message(fmt.Sprintf("Successfully created and joined group - %s", groupName))
}

func (s *Server) joinGroup(c *Client, groupName string) {
	group, ok := s.groups[groupName]
	if !ok {
		c.message("ERR: this group does not exist, if you want, you can create one with /create command!")
		return
	}

	c.group = group
	group.clients = append(group.clients, c)
	c.message(fmt.Sprintf("Successfully joined group - %s", groupName))
}
