package main

import (
	"log"
	"os"

	"github.com/Rachit-Gandhi/gator/internal/commands"
	"github.com/Rachit-Gandhi/gator/internal/config"
)

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	state := commands.State{
		Cfg: &config,
	}
	commandsMap := commands.Commands{
		Mux: make(map[string]func(s *commands.State, cmd commands.Command) error),
	}
	commandsMap.Register("login", commands.HandlerLogin)
	cliArgs := os.Args
	if len(cliArgs) < 2 {
		log.Fatal("Minimum of two arguments are expected.")
	}
	cmd := commands.Command{
		TriggerName: cliArgs[1],
		StringArgs:  cliArgs[2:],
	}
	if err := commandsMap.Run(&state, cmd); err != nil {
		log.Fatalf("Error running command: %v", err)
	}
}
