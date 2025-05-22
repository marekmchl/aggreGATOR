package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // imported for side effects

	"github.com/marekmchl/aggreGATOR/internal/cli"
	"github.com/marekmchl/aggreGATOR/internal/config"
	"github.com/marekmchl/aggreGATOR/internal/database"
	"github.com/marekmchl/aggreGATOR/internal/state"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("failed - %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Printf("failed - %v", err)
		os.Exit(1)
	}

	sta := &state.State{
		Config: &cfg,
	}

	dbQueries := database.New(db)
	sta.DB = dbQueries

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
