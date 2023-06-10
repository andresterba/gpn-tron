package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
)

type Player struct {
	ID         string
	MyPosition string
	X          int
	Y          int
}

type PositionData struct {
	Value    string
	PlayerID int
}

type Game struct {
	Map           [][]PositionData
	width         int
	height        int
	tickCounter   int
	DeadPlayerIDs []int
}

func buildGameMap(x, y int) [][]PositionData {
	gameMap := make([][]PositionData, y)
	for i := range gameMap {
		gameMap[i] = make([]PositionData, x)
		for j := range gameMap[i] {
			gameMap[i][j] = PositionData{
				Value:    FREE,
				PlayerID: 0,
			}
		}
	}

	return gameMap
}

func printGameMap(gameMap [][]PositionData) {
	for _, v := range gameMap {
		fmt.Printf("%+v\n", v)
	}
}

func move(direction string) {
	fmt.Printf("I'm going to move %s\n", direction)
	_, err := Write(conn, "move|"+direction+"\n")
	if err != nil {
		log.Printf("Write Error: %s\n", err)
	}
}

func moveRandom() {
	direction := rand.Intn(4)

	_, err := Write(conn, "move|"+directions[direction]+"\n")
	if err != nil {
		log.Printf("Write Error: %s\n", err)
	}
}

func isPositionFree(x, y int) bool {
	position := game.Map[x][y]
	if position.Value != OCCUPIED {
		return true
	}
	if contains(game.DeadPlayerIDs, position.PlayerID) {
		return true
	}

	return false
}

func Write(conn net.Conn, content string) (int, error) {
	writer := bufio.NewWriter(conn)
	number, err := writer.WriteString(content)
	if err == nil {
		err = writer.Flush()
	}

	return number, err
}
