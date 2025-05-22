package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/marekmchl/aggreGATOR/internal/database"
	"github.com/marekmchl/aggreGATOR/internal/state"
)

func handlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) <= 0 {
		return fmt.Errorf("command login requires a username")
	}
	s.Config.SetUser(cmd.Args[0])
	fmt.Println("success - the new user has been set")
	return nil
}

func handlerRegister(s *state.State, cmd Command) error {
	if len(cmd.Args) <= 0 {
		return fmt.Errorf("command register requires a name")
	}
	_, err := s.DB.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Args[0],
		},
	)
	if err != nil {
		return fmt.Errorf("creation of user %v was unsuccessful - %v", cmd.Args[0], err)
	}
	s.Config.SetUser(cmd.Args[0])

	fmt.Println("success - the new user has been registered")
	return nil
}
