package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	funk "github.com/thoas/go-funk"
)

// A Server manages the discord irc interface
type Server struct {
	name    string
	version string
	clients []*Client
}

// NewServer makes a new server with no clients
func NewServer(name string) Server {
	return Server{
		name:    name,
		version: "nosserver-1.1",
		clients: make([]*Client, 0),
	}
}

// Run starts the server loop
func (server *Server) Run() {

	listener, err := net.Listen("tcp", "localhost:6667")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on localhost:6667")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(server, conn)
	}
}

// TryRegister will try to register a valid client
func (server *Server) tryRegister(c *Client) {
	if c.registered {
		return
	}

	if c.name == "" || c.nick == "" {
		return
	}

	c.Send(server.name, RPLWelcome, c.nick, "Hi, welcome to IRC")
	c.Send(server.name, RPLYourhost, c.nick, fmt.Sprintf("Your host is %s, running version %s", server.name, server.version))
	c.Send(server.name, RPLCreated, c.nick, "This server was created sometime")
	// TODO: handle modes properly
	c.Send(server.name, RPLMyinfo, c.nick, server.name, server.version, "o o")

	// Motd
	c.Send(server.name, RPLMotdstart, c.nick, fmt.Sprintf("- %s Message of the day - ", server.name))
	c.Send(server.name, RPLMotd, c.nick, "- hello world!")
	c.Send(server.name, RPLEndofmotd, c.nick, "End of motd command")

	// luser
	c.Send(server.name, RPLLuserclient, c.nick, fmt.Sprintf("There are %d clients and 0 services on 1 server", len(server.clients)))
	c.registered = true
}

// Handles incoming requests.
func handleRequest(server *Server, conn net.Conn) {
	reader := bufio.NewReader(conn)
	client := NewClient(conn)
	server.clients = append(server.clients, &client)
	defer func() {
		conn.Close()
		// Delete client
		i := funk.IndexOf(server.clients, &client)
		server.clients = append(server.clients[:i], server.clients[i+1:]...)
	}()

	for {
		data, err := reader.ReadString('\n')
		data = strings.TrimSpace(data)
		if err != nil {
			break
		}
		fmt.Println("<<", data)
		split := strings.Split(data, " ")
		cmd := strings.ToUpper(split[0])
		args := split[1:]

		client.runCommand(server, cmd, args)
		// 	fmt.Fprintf(conn, ":%s 001 %s :Hi, welcome to IRC\r\n", servername, client.Nick)

		// 	fmt.Fprintf(conn, ":%[1]s 002 %[2]s :Your host is %[1]s, running version nosserver-1.1\r\n", servername, client.Nick)
		// 	fmt.Fprintf(conn, ":%s 003 %s :This server was created sometime\r\n", servername, client.Nick)
		// 	fmt.Fprintf(conn, ":%[1]s 004 %[2]s %[1]s nosserver-1.1 o o\r\n", servername, client.Nick)

		// 	fmt.Fprintf(conn, ":%[1]s 375 %[2]s := %[1]s Message of the day -\r\n", servername, client.Nick)
		// 	fmt.Fprintf(conn, ":%s 372 %s :- Hiya!\r\n", servername, client.Nick)
		// 	fmt.Fprintf(conn, ":%s 376 %s :End of /MOTD command\r\n", servername, client.Nick)

		// 	server.clients = append(server.clients, client)
		// 	fmt.Fprintf(conn, ":%s 251 %s :There are %d clients and 0 services on 1 server\r\n", servername, client.Nick, len(server.clients))
		// 	registered = true

		// 	// fmt.Fprintf(conn, ":noskcaj JOIN :#hack")
		// 	// // conn.Write([]byte(":noskcaj JOIN :#hack"))
		// 	// fmt.Fprintf(conn, ":nos 353 noskcaj = #hack :@noskcaj")
		// 	// // conn.Write([]byte(":nos 353 noskcaj = #hack :@noskcaj"))
		// 	// fmt.Fprintf(conn, ":nos 366 noskcaj #hack :End of /NAMES list.")
		// 	// // conn.Write([]byte(":nos 366 noskcaj #hack :End of /NAMES list."))
		// }
	}
}
