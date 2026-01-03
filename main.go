package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Rachit-Gandhi/gator/internal/commands"
	"github.com/Rachit-Gandhi/gator/internal/config"
	"github.com/Rachit-Gandhi/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	state := commands.State{
		Db:  dbQueries,
		Cfg: &cfg,
	}
	commandsMap := commands.Commands{
		Mux: make(map[string]func(s *commands.State, cmd commands.Command) error),
	}
	err = commandsMap.Register("login", commands.HandlerLogin)
	if err != nil {
		log.Fatalf("registration of command failed: %v", err)
	}
	err = commandsMap.Register("register", commands.Register)
	if err != nil {
		log.Fatalf("registration of command failed: %v", err)
	}
	err = commandsMap.Register("reset", commands.Reset)
	if err != nil {
		log.Fatalf("registration of command failed: %v", err)
	}
	err = commandsMap.Register("users", commands.GetUsers)
	if err != nil {
		log.Fatalf("registration of command failed: %v", err)
	}
	err = commandsMap.Register("agg", commands.Aggregate)
	if err != nil {
		log.Fatalf("registration of command failed: %v", err)
	}
	err = commandsMap.Register("addfeed", commands.AddFeed)
	if err != nil {
		log.Fatalf("registration of command failed: %v", err)
	}
	err = commandsMap.Register("feeds", commands.GetFeeds)
	if err != nil {
		log.Fatalf("registration of command failed: %v", err)
	}
	err = commandsMap.Register("follow", commands.CreateFeedFollow)
	if err != nil {
		log.Fatalf("registration of command failed: %v", err)
	}
	err = commandsMap.Register("following", commands.GetFeedFollowsForUser)
	if err != nil {
		log.Fatalf("registration of command failed: %v", err)
	}
	cliArgs := os.Args
	if len(cliArgs) < 2 {
		log.Fatal("expected a command name")
	}
	cmd := commands.Command{
		TriggerName: cliArgs[1],
		StringArgs:  cliArgs[2:],
	}
	if err := commandsMap.Run(&state, cmd); err != nil {
		log.Fatalf("error running command: %v", err)
	}
}
