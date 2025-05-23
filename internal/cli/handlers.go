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

func handlerAddfeed(s *state.State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("addfeed needs both the feed name and url")
	}
	current_user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't add feed - %v", err)
	}
	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    current_user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't add feed - %v", err)
	}

	if err := handlerFollow(s, Command{
		Name: "follow",
		Args: []string{feed.Url},
	}); err != nil {
		return fmt.Errorf("couldn't create a follow - %v", err)
	}

	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state.State, cmd Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("getting feeds unsuccessful - %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("%v | %v | %v\n", feed.FeedName, feed.Url, feed.UserName)
	}
	return nil
}

func handlerFollow(s *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("command follow needs url")
	}

	user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("getting user info unsuccessful - %v", err)
	}

	feed, err := s.DB.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("getting feed info unsuccessful - %v", err)
	}

	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("creating follow was unsuccessful - %v", err)
	}

	fmt.Printf("%v | %v", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerFollowing(s *state.State, cmd Command) error {
	user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("getting user info was unsuccessful - %v", err)
	}
	feeds, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("getting user's feeds was unsuccessful - %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("%v\n", feed.FeedName)
	}
	return nil
}
