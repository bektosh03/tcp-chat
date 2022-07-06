package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Client struct {
	username string
	conn     net.Conn
	group    *Group
}

func NewClient(conn net.Conn) (Client, error) {
	c := Client{
		username: "",
		conn:     conn,
	}

	if err := c.getUsername(); err != nil {
		return Client{}, err
	}

	return c, nil
}

func (c *Client) getUsername() error {
	if err := c.message("Enter your name: "); err != nil {
		return err
	}

	input, err := c.readInput()
	if err != nil {
		return err
	}

	c.username = input

	return nil
}

func (c *Client) readInput() (string, error) {
	msg, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		log.Printf("error while reading from conn: %v", err)
		return "", err
	}
	if errors.Is(err, io.EOF) {
		msg = "/exit"
	}

	return strings.Trim(msg, "\n"), nil
}

func (c *Client) message(msg string) error {
	if _, err := c.conn.Write([]byte(fmt.Sprintf("> %s\n", msg))); err != nil {
		return err
	}

	return nil
}
