package timesender

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func StartClient() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 256)
	for {
		_, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		_, err = io.WriteString(os.Stdout, fmt.Sprintf("Custom output! %s", string(buf)))
		if err != nil {
			log.Println(err)
		}
	}
}
