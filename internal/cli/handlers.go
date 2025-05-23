package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/marekmchl/aggreGATOR/internal/database"
	"github.com/marekmchl/aggreGATOR/internal/rss"
	"github.com/marekmchl/aggreGATOR/internal/state"
)

func handlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) <= 0 {
		return fmt.Errorf("command login requires a username")
	}
	_, err := s.DB.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("user not found")
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

func handlerReset(s *state.State, cmd Command) error {
	if err := s.DB.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("resetting users was unsuccessful - %v", err)
	}
	return nil
}

func handlerUsers(s *state.State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("getting users unsuccessful - %v", err)
	}
	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("%v (current)\n", user.Name)
		} else {
			fmt.Println(user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state.State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("aggregation unsuccessful - %v", err)
	}
	fmt.Println(feed)
	return nil
}
