package main

import (
	"fmt"
	"math/rand"
)

func main() {

	fmt.Print("Введите интервал случайного числа: ")
	var a, b int
	fmt.Scan(&a, &b)

	var randomNum int

	randomNum = rand.Intn(b-a) + a

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
