package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	//SERVER_ADDR = "gpn-tron.duckdns.org:4000"
	SERVER_ADDR      = "localhost:4000"
	DELIMITER   byte = '\n'
	OCCUPIED         = "X"
	FREE             = "O"
)

var (
	conn       net.Conn
	game       Game
	myPlayer   Player
	directions = []string{"up", "down", "left", "right"}

	Username string
	Password string
)

func main() {
	game = Game{}
	game.tickCounter = 0
	game.DeadPlayerIDs = []int{}

	Username = getEnviromentVariableOrDefaultAsString("BOT_USERNAME", "test-user")
	Password = getEnviromentVariableOrDefaultAsString("PASSWORD", "Hallo123")

	myPlayer = Player{}

	connection, err := net.Dial("tcp", SERVER_ADDR)
	if err != nil {
		panic(err)
	}
	conn = connection
	defer connection.Close()

	tronHandlers := map[string]protocolHandler{
		"motd":    handleMotd,
		"error":   handleError,
		"game":    handleGame,
		"pos":     handlePos,
		"tick\n":  handleTick,
		"die":     handleDie,
		"lose":    handleLose,
		"win":     handleWin,
		"message": handleMessage,
		"player":  handlePlayerMessage,
	}

	join(Username, Password)

	reader := bufio.NewReader(connection)
	for {
		line, err := reader.ReadBytes(byte('\n'))
		switch err {
		case nil:
			break
		case io.EOF:
			fmt.Println("EOF")
			os.Exit(1)
		default:
			fmt.Println("ERROR", err)
			os.Exit(1)
		}

		determineType(tronHandlers, string(line))
	}
}

func determineWhereToMove() {
	// is up possible
	upX := myPlayer.X
	upY := (myPlayer.Y - 1 + game.height) % game.height
	if isPositionFree(upX, upY) {
		fmt.Printf("Should be free: x=%d y=%d\n", upX, upY)
		move("up")
		return
	}

	// is left possible
	leftX := (myPlayer.X - 1 + game.width) % game.width
	leftY := myPlayer.Y
	if isPositionFree(leftX, leftY) {
		fmt.Printf("Should be free: x=%d y=%d\n", leftX, leftY)
		move("left")
		return
	}

	// is right possible
	rightX := (myPlayer.X + 1) % game.width
	rightY := myPlayer.Y
	if isPositionFree(rightX, rightY) {
		fmt.Printf("Should be free: x=%d y=%d\n", rightX, rightY)
		move("right")
		return
	}

	// is down possible
	downX := myPlayer.X
	downY := (myPlayer.Y + 1) % game.height

	if isPositionFree(downX, downY) {
		fmt.Printf("Should be free: x=%d y=%d\n", downX, downY)
		move("down")
		return
	}

	fmt.Println("########################################")
	fmt.Println("If you can see this we are going to die!")
	fmt.Println("########################################")
}
