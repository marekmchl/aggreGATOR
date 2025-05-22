package cli

import (
	"fmt"

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
