package main

// NickCommand handles the nick command
func nickCommand(server *Server, c *Client, args []string) {
	c.nick = args[0]
}

// UserCommand handles the user command
func userCommand(server *Server, c *Client, args []string) {
	c.name = args[0]
	c.realname = args[3]
}

func pingCommand(server *Server, c *Client, args []string) {
	c.Send(server.name, "PONG")
}
