package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/mbrostami/gcron/cron"
)

// ListenUDP Start listening on udp port
func ListenUDP(host string, port string) {
	// Listen for incoming connections.
	conn, err := net.ListenPacket("udp", host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer conn.Close()
	fmt.Println("Listening on udp: " + host + ":" + port)
	for {
		// Listen for an incoming connection.
		// Handle connections in a new goroutine.
		handleUDPRequest(conn)
	}
}

// Handles incoming requests.
func handleUDPRequest(conn net.PacketConn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, raddr, err := conn.ReadFrom(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	tmpbuff := bytes.NewBuffer(buf)
	tmpstruct := new(cron.Task)
	gobobj := gob.NewDecoder(tmpbuff)
	gobobj.Decode(tmpstruct)
	fmt.Printf("%+v", string(tmpstruct.Output))
	fmt.Printf("%+v", tmpstruct)
	// Send a response back to person contacting us.
	conn.WriteTo([]byte("Message received."), raddr)
	// Close the connection when you're done with it.
}
