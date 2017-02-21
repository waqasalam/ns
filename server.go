package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"netns"
	"os"
	"runtime"
	//	"syscall"
	"time"
	"utils/netUtils"
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
func send_handler(sendc chan *Message, done chan struct{}) {

	for msg := range sendc {
		msg := msg
		fmt.Printf("message rcvd %s on", msg.name)
		msg.conn.Write([]byte("Message received." + msg.name))
		msg.conn.Close()
		done <- struct{}{}
	}
}

func connect(path string) {
	runtime.LockOSThread()
	fmt.Println("Do connect", path)

	socket, err := netUtils.ConnectSocket("tcp", "0.0.0.0:3333", "0.0.0.0:179")
	defer netUtils.CloseSocket(socket)
	if err != nil {
		log.Fatal("Socket connect failed")
	}

	err = netUtils.Connect(socket, "tcp", "0.0.0.0:3333", "0.0.0.0:179", time.Duration(5)*time.Second)
	if err != nil {
		log.Fatal("Socket connect failed")
	}

	//	conn, err := netUtils.ConvertFdToConn(socket)
	//	if err != nil {
	//		fmt.Println("conn didn't go thru")
	//	}

}

func main() {

	fmt.Println("Main program started", unix.Gettid)

	done := make(chan struct{})
	sendc := make(chan *Message)
	go send_handler(sendc, done)

	nsfunc := func(path string) {
		runtime.LockOSThread()
		fmt.Println(path)
		if path != "" {
			newns, err := netns.GetFromPath(path)
			if err != nil {
				log.Fatal(err)
			}
			if err := netns.SetNs(newns); err != nil {
				log.Fatal(err)
			}
		}
		l, err := net.Listen("tcp4", "0.0.0.0:3333")
		if err != nil {
			fmt.Println("Error listening:", err.Error(), path)
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
				sendc <- msg
			}
			fmt.Println("Done", path)
			l.Close()
		}
		// Do accept in a separate go routine
		go accept_func()
	}

	go nsfunc("coke")
	go nsfunc("waqas")
	go nsfunc("")
	for i := 0; i < 6; i++ {
		<-done
	}
	//Start listening on the namespace again check if it works.

	go nsfunc("coke")
	go nsfunc("waqas")
	go nsfunc("")
	for i := 0; i < 4; i++ {
		<-done
	}

	close(sendc)

	//	ns, err := GetFromThread()
	//	if err != nil {
	//		t.Fatal(err)
	//	}

}
