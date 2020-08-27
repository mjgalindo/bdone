package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/gen2brain/beeep"
)

func serve(port string) (err error) {
	pc, err := net.ListenPacket("udp", "0.0.0.0:"+port)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	for {

		buffer := make([]byte, 1024)

		_, _, err = pc.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}

		beeep.Beep(beeep.DefaultFreq, 500)

		fmt.Println("Message received = " + string(buffer))
		err := beeep.Notify(string(buffer), string(buffer), "")
		if err != nil {
			panic(err)
		}
	}
}

// a message and waits for a response coming back from the server
// that it initially targetted.
func client(address string, port string, message string) (err error) {
	raddr, err := net.ResolveUDPAddr("udp", address+":"+port)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return
	}

	// Closes the underlying file descriptor associated with the,
	// socket so that it no longer refers to any file.
	defer conn.Close()

	doneChan := make(chan error, 1)

	rdr := strings.NewReader(message)
	_, err = io.Copy(conn, rdr)
	if err != nil {
		doneChan <- err
		return
	}

	return
}

func main() {
	isServer := flag.Bool("s", false, "Serve instead of send")
	addr := flag.String("a", "localhost", "The address to send the msg")
	port := flag.String("p", "25651", "The port")
	msg := flag.String("m", "Hello", "The message to send")
	flag.Parse()
	if *isServer {
		serve(*port)
	} else {
		client(*addr, *port, *msg)
	}
}
