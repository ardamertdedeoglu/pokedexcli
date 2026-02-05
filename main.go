package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "Pokedex > ",
		HistoryFile:     "pokedex_history.txt", // Komutlar bu dosyaya kaydedilir
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		command, err := rl.Readline()
		if err != nil { // Ctrl+C veya Ctrl+D (EOF) durumunda döngüden çıkar
			break
		}

		command = strings.TrimSpace(command)
		if command == "" {
			continue
		}
		clean_command := cleanInput(command)
		real_command := clean_command[0]
		res, ok := commands[real_command]
		if !ok {
			fmt.Println("Unknown command")
		}
		args := clean_command[1:]
		err = res.callback(&res.config, args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}
