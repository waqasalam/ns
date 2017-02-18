package netns

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"testing"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func TestNs(t *testing.T) {
	//origns, err := GetFromThread()
	//if err != nil {
	//	t.Fatal(err)
	//}
	done := make(chan struct{})
	nsfunc := func(path string) {
		runtime.LockOSThread()
		fmt.Println(path)
		newns, err := GetFromPath(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := SetNs(newns); err != nil {
			t.Fatal(err)
		}

		l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
		if err != nil {
			fmt.Println("Error listening:", err.Error())
			os.Exit(1)
		}
		fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
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

		fmt.Println("msg received")
	}
	//	ns, err := GetFromThread()
	//	if err != nil {
	//		t.Fatal(err)
	//	}

}
