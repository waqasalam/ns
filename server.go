package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"netns"
	"os"
	"runtime"
)

const (
	//CONN_HOST = ""
	CONN_PORT = "3333"
	//CONN_TYPE = "tcp"
)

type Message struct {
	conn net.Conn
	name string
}

// This is doing send message for all the vrfs
func send_handler(c chan *Message, done chan struct{}) {

	for msg := range c {
		msg := msg
		fmt.Printf("message rcvd %s on", msg.name)
		msg.conn.Write([]byte("Message received." + msg.name))
		msg.conn.Close()
		done <- struct{}{}
	}
}

func main() {

	fmt.Println("Main program started", unix.Gettid)

	done := make(chan struct{})
	msgc := make(chan *Message)
	go handler(msgc, done)

	nsfunc := func(path string) {
		runtime.LockOSThread()
		fmt.Println(path)

		newns, err := netns.GetFromPath(path)
		if err != nil {
			log.Fatal(err)
		}
		if err := netns.SetNs(newns); err != nil {
			log.Fatal(err)
		}

		l, err := net.Listen("tcp4", "0.0.0.0:3333")
		if err != nil {
			fmt.Println("Error listening:", err.Error())
			os.Exit(1)
		}

		accept_func := func() {
			fmt.Println("Listening on "+":"+CONN_PORT, unix.Gettid)
			for i := 0; i < 2; i++ {
				// Listen for an incoming connection.
				conn, err := l.Accept()
				if err != nil {
					fmt.Println("Error accepting: ", err.Error())
					os.Exit(1)
				}

				buf := make([]byte, 1024)
				_, err = conn.Read(buf)
				if err != nil {
					fmt.Println("Error Reading", err.Error())
				}
				msg := &Message{conn: conn, name: path}
				msgc <- msg
			}
			fmt.Println("Done", path)
			l.Close()
		}
		// Do accept in a separate go routine
		go accept_func()
	}

	go nsfunc("coke")
	go nsfunc("waqas")

	for i := 0; i < 4; i++ {
		<-done
	}
	//Start listening on the namespace again check if it works.

	go nsfunc("coke")
	go nsfunc("waqas")

	for i := 0; i < 4; i++ {
		<-done
	}

	close(msgc)

	//	ns, err := GetFromThread()
	//	if err != nil {
	//		t.Fatal(err)
	//	}

}
