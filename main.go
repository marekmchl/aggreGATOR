package main

import (
	"fmt"

	"github.com/marekmchl/aggreGATOR/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	cfg.SetUser("marekmchl")
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("db_url: %v\ncurrent_user_name: %v", cfg.DbURL, cfg.CurrentUserName)
}
