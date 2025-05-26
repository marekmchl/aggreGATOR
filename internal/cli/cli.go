package cli

import (
	"fmt"

	"github.com/marekmchl/aggreGATOR/internal/state"
)

type Command struct {
	Name string
	Args []string
}

func GetCommands() commands {
	cmds := commands{
		commandMap: make(map[string]func(*state.State, Command) error, 1),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	return cmds
}

type commands struct {
	commandMap map[string]func(*state.State, Command) error
}

func (c *commands) Run(s *state.State, cmd Command) error {
	command, found := c.commandMap[cmd.Name]
	if !found {
		return fmt.Errorf("command not found")
	}

	if err := command(s, cmd); err != nil {
		return err
	}

	return nil
}
func (c *commands) register(name string, f func(*state.State, Command) error) {
	c.commandMap[name] = f
}
