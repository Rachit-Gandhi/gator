package commands

import (
	"fmt"

	"github.com/Rachit-Gandhi/gator/internal/config"
)

type State struct {
	Cfg *config.Config
}

type Command struct {
	TriggerName string
	StringArgs  []string
}

type Commands struct {
	Mux map[string]func(s *State, cmd Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	v, ok := c.Mux[cmd.TriggerName]
	if !ok {
		return fmt.Errorf("command not found")
	}
	err := v(s, cmd)
	if err != nil {
		return fmt.Errorf("error executing error: %w", err)
	}
	return nil
}

func (c *Commands) Register(TriggerName string, f func(*State, Command) error) error {
	_, ok := c.Mux[TriggerName]
	if ok {
		return fmt.Errorf("command trigger %v already registerd, cannot register again", TriggerName)
	}

	c.Mux[TriggerName] = f
	return nil
}
