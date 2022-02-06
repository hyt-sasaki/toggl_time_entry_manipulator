package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jason0x43/go-alfred"
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

    firstCall := os.Args[2] == ""
    commands, err := initializeCommands(workflow, firstCall)
    if err != nil {
        log.Fatalln(err)
        os.Exit(1)
    }
    workflow.Run(commands)
}

