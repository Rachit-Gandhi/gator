package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/Rachit-Gandhi/gator/internal/database"
	"github.com/Rachit-Gandhi/gator/internal/rss"
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

func GetUser(s *State, cmd Command) (database.User, error) {
	if len(cmd.StringArgs) != 1 {
		return database.User{}, fmt.Errorf("the user handler needs one argument username of type string")
	}
	user, err := s.Db.GetUser(context.Background(), cmd.StringArgs[0])
	if err != nil {
		return database.User{}, fmt.Errorf("error getting user with this username: %w", err)
	}
	return user, nil
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

func Aggregate(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 0 {
		return fmt.Errorf("the users handler doesn't accept cli argument")
	}
	r, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("the xml fetch failed: %w", err)
	}
	fmt.Println(r)
	return nil

}

func AddFeed(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 2 {
		return fmt.Errorf(("add feed handlers expects 2 arguments: name & url"))
	}
	id := uuid.New()
	currUserName := s.Cfg.GetUser()
	currUser, err := s.Db.GetUser(context.Background(), currUserName)
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}
	user_id := currUser.ID
	name, url := cmd.StringArgs[0], cmd.StringArgs[1]
	newFeed := database.AddFeedParams{
		ID:     id,
		UserID: user_id,
		Name:   name,
		Url:    url,
	}
	_, err = s.Db.AddFeed(context.Background(), newFeed)
	fmt.Printf("feed aded: %v\n", newFeed)
	return nil
}

func GetFeeds(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 0 {
		return fmt.Errorf("get feeds doesn't accept any arguments")
	}
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}
	for _, feed := range feeds {
		userName, err := s.Db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting username: %w", err)
		}
		fmt.Printf("Name: %v, URL: %v, userName: %v\n", feed.Name, feed.Url, userName.Name)
	}
	return nil
}
