package main

import (
	"time"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
)

type Config struct {
	APIKey           string `desc:"Toggl API key"`
}
type Cache struct {
	Workspace int
	Account   toggl.Account
	Time      time.Time
}


func checkRefresh() error {
	if time.Now().Sub(cache.Time).Minutes() < 5.0 {
		return nil
	}

	dlog.Println("Refreshing cache...")
	err := refresh()
	if err != nil {
		dlog.Println("Error refreshing cache:", err)
	}
	return err
}

func refresh() error {
	s := toggl.OpenSession(config.APIKey)
	account, err := s.GetAccount()
	if err != nil {
		return err
	}

	dlog.Printf("got account: %#v", account)

	cache.Time = time.Now()
	cache.Account = account
	cache.Workspace = account.Data.Workspaces[0].ID
	return alfred.SaveJSON(cacheFile, &cache)
}
