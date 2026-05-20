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
			continue
		}
		command := formatted[0]
		pokeCommands := getCommands()
		cmd, ok := pokeCommands[command]
		if !ok {
			fmt.Println("Unkown command")
		} else {
			err := cmd.callback()
			if err != nil {
				logError := fmt.Sprintf("'%s' error: %v", command, err)
				log.Println(logError)
			}
		}
	}
}
