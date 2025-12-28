package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/Rachit-Gandhi/gator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	_, ok := s.Db.GetUser(context.Background(), cmd.StringArgs[0])
	if ok != nil {
		return fmt.Errorf("user doesn't exist in the database")
	}
	err := s.Cfg.SetUser(cmd.StringArgs[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set as: %v\n", cmd.StringArgs[0])
	return nil
}

func Register(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	_, ok := s.Db.GetUser(context.Background(), cmd.StringArgs[0])
	if ok == nil {
		return fmt.Errorf("user with that username already exists")
	}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.StringArgs[0],
	}
	_, err := s.Db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("user creation in db failed, %w", err)
	}
	err = s.Cfg.SetUser(cmd.StringArgs[0])
	if err != nil {
		return err
	}
	fmt.Printf("new user has been created as: %v\n", cmd.StringArgs[0])
	return nil
}

func Reset(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 0 {
		return fmt.Errorf("the reset handler doesn't accept cli argument")
	}
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting all users: %w\n", err)
	}
	fmt.Println("all users are reset")
	return nil
}

func GetUsers(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 0 {
		return fmt.Errorf("the users handler doesn't accept cli argument")
	}
	usrs, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting all users: %w\n", err)
	}
	currentUser := s.Cfg.GetUser()
	for _, usr := range usrs {
		fmt.Print("* ")
		if usr.Name != currentUser {
			fmt.Printf("%v\n", usr.Name)
		} else {
			fmt.Printf("%v (current)\n", currentUser)
		}
	}
	return nil
}
