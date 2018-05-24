package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// A Client that has joined the server
type Client struct {
	name     string
	nick     string
	realname string

	registered bool

	conn net.Conn
}

// NewClient creates a new client with a connection
func NewClient(conn net.Conn) Client {
	return Client{
		registered: false,
		conn:       conn,
	}
}

// RunCommand runs a command
func (c *Client) runCommand(server *Server, cmd string, args []string) {
	switch cmd {
	case "NICK":
		nickCommand(server, c, args)
	case "USER":
		userCommand(server, c, args)
	case "PING":
		pingCommand(server, c, args)
	default:
		fmt.Printf("Unknown command \"%s\" %s\n", cmd, args)
	}

	if !c.registered {
		server.tryRegister(c)
	}
}

// Send sends a formatted message to a client
func (c *Client) Send(prefix, command string, params ...string) error {
	var buf string

	if len(prefix) > 0 {
		buf += fmt.Sprintf(":%s ", prefix)
	}

	buf += command

	for i, param := range params {
		buf += " "
		if len(param) < 1 || strings.Contains(param, " ") || param[0] == ':' {
			if i != len(params)-1 {
				return errors.New("irc: Cannot have an empty param, a param with spaces, or a param that starts with ':' before the last parameter")
			}
			buf += ":"
		}
		buf += param
	}
	buf += "\r\n"

	fmt.Fprint(c.conn, buf)
	return nil
}
