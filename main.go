package main

import (
	"fmt"
	"os"

	"github.com/marekmchl/aggreGATOR/internal/cli"
	"github.com/marekmchl/aggreGATOR/internal/config"
	"github.com/marekmchl/aggreGATOR/internal/state"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("failed - %v\n", err)
		os.Exit(1)
	}

	sta := &state.State{
		Config: &cfg,
	}

	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Println("failed - no command found")
		os.Exit(1)
	}

	cmd := cli.Command{
		Name: args[0],
		Args: args[1:],
	}

	cmds := cli.GetCommands()

	if err := cmds.Run(sta, cmd); err != nil {
		fmt.Printf("failed - %v\n", err)
		os.Exit(1)
	}
}
