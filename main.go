package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	prompt := "Pokedex > "
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(prompt)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatalf("input error: %v", err)
		}
		input := scanner.Text()
		formatted := cleanInput(input)
		if len(formatted) == 0 {
			fmt.Println("no input received")
		} else {
			output := fmt.Sprintf("Your command was: %s", formatted[0])
			fmt.Println(output)
		}
	}
}
