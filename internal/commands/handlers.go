package commands

import "fmt"

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.StringArgs) == 0 || len(cmd.StringArgs) > 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	err := s.Cfg.SetUser(cmd.StringArgs[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set as: %v\n", cmd.StringArgs[0])
	return nil
}
