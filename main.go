package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jason0x43/go-alfred"
    "toggl_time_entry_manipulator/command/add"
    "toggl_time_entry_manipulator/command/list"
    "toggl_time_entry_manipulator/command/get"
    "toggl_time_entry_manipulator/command/stop"
    "toggl_time_entry_manipulator/command/modify"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator]", log.LstdFlags)

func main() {
    var workflow alfred.Workflow
    var err error
    if workflow, err = alfred.OpenWorkflow("..", true); err != nil {
        fmt.Printf("Error: %s", err)
        os.Exit(1)
    }

    dlog.Printf("cache file dir: %s", workflow.CacheDir())
    dlog.Printf("config file dir: %s", workflow.DataDir())

    // firestore
    repo, err := initializeRepository(workflow)
    if err != nil {
        log.Fatalln(err)
        os.Exit(1)
    }

    workflow.Run([]alfred.Command{
        add.AddEntryCommand{
            Repo: repo,
        },
        list.ListEntryCommand{
            Repo: repo,
        },
        get.GetEntryCommand{
            Repo: repo,
        },
        stop.StopEntryCommand{
            Repo: repo,
        },
        modify.ModifyEntryCommand{
            Repo: repo,
        },
    })
}

