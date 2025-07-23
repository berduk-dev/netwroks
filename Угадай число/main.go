package main

import (
	"fmt"
	"math/rand"
)

func main() {

	randomNum := rand.Intn(10)

	var guessableNum int
	for {
		fmt.Scan(&guessableNum)
		if randomNum == guessableNum {
			fmt.Printf("Вы угадали число! %d\n", randomNum)
			break
		} else if randomNum > guessableNum {
			fmt.Println("Попробуй число побольше")
		} else {
			fmt.Println("Попробуй число поменьше")
		}
	}
}
