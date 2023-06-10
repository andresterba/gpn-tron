package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type (
	protocolHandler func(string)
)

func determineType(handlers map[string]protocolHandler, data string) {
	splittedData := strings.Split(data, "|")
	protocolType := splittedData[0]

	dataInput := strings.Join(splittedData[1:], "|")
	safeString := strings.ReplaceAll(dataInput, "\n", "")

	handlers[protocolType](safeString)
}

func handleMotd(data string) {
	fmt.Printf("TYPE: motd\t PAYLOAD: %s\n", data)
}

func handleError(data string) {
	fmt.Printf("TYPE: error\t PAYLOAD: %s\n", data)
}

func login(username, password string) {
	loginContent := fmt.Sprintf("join|%s|%s\n", username, password)

	_, err := Write(conn, loginContent)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
	}

	log.Printf("Logged in as %s\n", username)
}

func handleGame(data string) {
	// game|100|100|5
	splittedData := strings.Split(data, "|")
	mapWidth := splittedData[0]
	mapHeight := splittedData[1]
	playerID := splittedData[2]

	myPlayer.ID = playerID

	mapWidthAsInt, err := strconv.Atoi(mapWidth)
	if err != nil {
		fmt.Println("Error during conversion")
		return
	}
	mapHeightAsInt, err := strconv.Atoi(mapHeight)
	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	game.width = mapWidthAsInt
	game.height = mapHeightAsInt
	game.Map = buildGameMap(game.width, game.height)

	fmt.Printf("Map Size %d x %d", game.height, game.width)
	fmt.Printf("My Player ID %s", myPlayer.ID)
	printGameMap(game.Map)
}

func handlePos(data string) {
	// pos|5|3|8
	splittedData := strings.Split(data, "|")

	playerID := splittedData[0]
	xPosition := splittedData[1]
	yPosition := splittedData[2]

	x := stringToInt(xPosition)
	y := stringToInt(yPosition)

	pos := fmt.Sprintf("%s-%s", xPosition, yPosition)

	if playerID == myPlayer.ID {
		myPlayer.MyPosition = pos
		myPlayer.X = stringToInt(xPosition)
		myPlayer.Y = stringToInt(yPosition)

		game.Map[x][y] = PositionData{
			Value:    OCCUPIED,
			PlayerID: stringToInt(playerID),
		}

		return
	}

	game.Map[x][y] = PositionData{
		Value:    OCCUPIED,
		PlayerID: stringToInt(playerID),
	}
}

func handleTick(data string) {
	game.tickCounter++

	fmt.Printf("Tick is %d\n", game.tickCounter)
	fmt.Printf("Current Player Position x=%d y=%d\n", myPlayer.X, myPlayer.Y)

	determineWhereToMove()
}

func handleDie(data string) {
	// Example (4 dead player): die|5|8|9|13

	splittedData := strings.Split(data, "|")

	for _, v := range splittedData {
		fmt.Printf("Player %s died! \n", v)
		playerID := stringToInt(v)
		game.DeadPlayerIDs = append(game.DeadPlayerIDs, playerID)
	}

}

func handleMessage(data string) {
	// message|5|I am so cool
	//fmt.Printf("TYPE: message\t PAYLOAD: %s\n", data)

	//splittedData := strings.Split(data, "|")
	//playerID := splittedData[0]
	//message := splittedData[1]

	//positionData := fmt.Sprintf("playerID: %s message: %s", playerID, message)
	//fmt.Printf("TYPE: message\t PAYLOAD: %s\n", positionData)
}

func handleLose(data string) {
	// message|5|I am so cool
	//fmt.Printf("TYPE: lose\t PAYLOAD: %s\n", data)

	splittedData := strings.Split(data, "|")
	amountWin := splittedData[0]
	amountLose := splittedData[1]

	positionData := fmt.Sprintf("wins: %s lose: %s", amountWin, amountLose)

	fmt.Printf("TYPE: lose\t PAYLOAD: %s\n", positionData)

	// Reset Game
	game = Game{}
	game.tickCounter = 0
	game.DeadPlayerIDs = []int{}

	myPlayer = Player{}
}

func handleWin(data string) {
	// message|5|I am so cool
	//fmt.Printf("TYPE: win\t PAYLOAD: %s\n", data)

	splittedData := strings.Split(data, "|")
	amountWin := splittedData[0]
	amountLose := splittedData[1]

	positionData := fmt.Sprintf("win: %s lose: %s", amountWin, amountLose)

	fmt.Printf("TYPE: win\t PAYLOAD: %s\n", positionData)
}
