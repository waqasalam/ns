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

func main() {
	//origns, err := GetFromThread()
	//if err != nil {
	//	t.Fatal(err)
	//}
	done := make(chan struct{})
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
			conn.Write([]byte("Message received."))
			conn.Close()
			fmt.Println("Message on ", path)
			done <- struct{}{}
		}
		fmt.Println("Done", path)
	}
	go nsfunc("coke")
	go nsfunc("waqas")

	for i := 0; i < 4; i++ {
		<-done
	}
	//	ns, err := GetFromThread()
	//	if err != nil {
	//		t.Fatal(err)
	//	}

}
