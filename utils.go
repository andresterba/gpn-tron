package main

import (
	"os"
	"strconv"
)

func stringToInt(input string) int {
	numberAsInt, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}

	return numberAsInt
}

func contains(s []int, value int) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}

	return false
}

func getEnviromentVariableOrDefaultAsString(name string, def string) string {
	env := os.Getenv(name)
	if env == "" {
		return def
	}

	return env
}
