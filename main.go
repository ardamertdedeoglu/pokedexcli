package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		available := scanner.Scan()
		if !available {
			fmt.Println("error getting new token")
			return
		}
		command := scanner.Text()
		clean_command := cleanInput(command)
		real_command := clean_command[0]
		res, ok := commands[real_command]
		if !ok {
			fmt.Println("Unknown command")
		}
		args := clean_command[1:]
		err := res.callback(&res.config, args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}
