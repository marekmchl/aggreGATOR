package cli

import (
	"context"
	"fmt"

	"github.com/marekmchl/aggreGATOR/internal/database"
	"github.com/marekmchl/aggreGATOR/internal/state"
)

func middlewareLoggedIn(handler func(s *state.State, cmd Command, user database.User) error) func(*state.State, Command) error {
	return func(s *state.State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("getting user info was unsuccessful - %v", err)
		}
		return handler(s, cmd, user)
	}
}
