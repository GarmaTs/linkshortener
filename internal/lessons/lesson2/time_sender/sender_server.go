package timesender

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var chMsg = make(chan string) // Канал для сообщений с консоли сервера
func StartListening(port string) {
	if len(port) == 0 {
		port = ":8000"
	}
	listener, err := net.Listen("tcp", "localhost"+port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

// Решение для первого пункта второго урока
func addToChMsg() {
	reader := bufio.NewReader(os.Stdin)
	msg, err := reader.ReadString('\n')
	if err != nil {
		log.Println("error in addToChMsg", err)
	}

	chMsg <- msg
}

func handleConn(c net.Conn) {
	defer c.Close()

	for {
		go addToChMsg()

		select {
		case msg := <-chMsg:
			_, err := io.WriteString(c, msg)
			if err != nil {
				return
			}
		default:
			_, err := io.WriteString(c, time.Now().Format("15:04:05\n\r"))
			if err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}
