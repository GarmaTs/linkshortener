package chat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

const userNameSplitter string = "_> "

type Game struct {
	numbers    [2]int
	operations [3]string
	answer     int
	equation   string
}

func (g *Game) newGame() {
	g.operations = [3]string{"+", "-", "*"}
}

func (g *Game) newEquation() (string, int) {
	g.numbers[0] = rand.Intn(10)
	g.numbers[1] = rand.Intn(10)
	operation := g.operations[rand.Intn(len(g.operations))]

	if g.numbers[0] == 0 {
		g.numbers[0] += 1
	}
	if g.numbers[1] == 0 {
		g.numbers[1] += 1
	}
	if operation == "-" && g.numbers[0] < g.numbers[1] {
		g.numbers[0], g.numbers[1] = g.numbers[1], g.numbers[0]
	}

	switch operation {
	case "+":
		g.answer = g.numbers[0] + g.numbers[1]
	case "-":
		g.answer = g.numbers[0] - g.numbers[1]
	case "*":
		g.answer = g.numbers[0] * g.numbers[1]
	}

	g.equation = fmt.Sprintf("%d %s %d", g.numbers[0], operation, g.numbers[1])
	return g.equation, g.answer
}

func (g Game) checkAnswer(num int) bool {
	if num != g.answer {
		return false
	}
	return true
}

func StartChatServer(port string) {
	if len(port) == 0 {
		port = ":8000"
	}

	listener, err := net.Listen("tcp", "localhost"+port)
	if err != nil {
		log.Fatal(err)
	}

	var game Game
	game.newGame()
	game.newEquation()

	go broadcaster(&game)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, &game)
	}
}

func broadcaster(game *Game) {
	clients := make(map[client]bool)
	for {
		select {
		case srcMsg := <-messages:
			for cli := range clients {
				// Рассылаем присланное сообщение
				cli <- srcMsg
			}

			arr := strings.Split(srcMsg, userNameSplitter)
			if len(arr) == 2 {
				// Проверяем ответ, и рассылаем сообщение о результате проверки
				var msg string
				num, err := strconv.Atoi(arr[1])
				if err != nil {
					continue
				}
				if game.checkAnswer(num) {
					game.newEquation()
					msg = "\r\n" + arr[0] + " gave right answer. New equation is " + game.equation
				} else {
					msg = "\r\n" + arr[1] + " is wrong answer. Try again"
				}

				for cli := range clients {
					cli <- msg
				}
			}

		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn, game *Game) {
	chLocal := make(chan string)
	go clientWriter(conn, chLocal)

	myReader := strings.NewReader("Welcome, input your name: ")
	io.Copy(conn, myReader)

	who := "Anonim"
	buf := bufio.NewScanner(conn)
	for buf.Scan() {
		who = buf.Text()
		break
	}

	chLocal <- fmt.Sprintf("Hello %s! What is answer for %s", who, game.equation)
	messages <- who + " has arrived"
	entering <- chLocal

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + userNameSplitter + input.Text()
	}
	leaving <- chLocal
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
