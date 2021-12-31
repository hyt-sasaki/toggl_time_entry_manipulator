package main

import (
	"fmt"
	"log"
	"os"
	"path"
    "google.golang.org/api/option"


	"github.com/jason0x43/go-alfred"
    "toggl_time_entry_manipulator/estimation_client"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator]", log.LstdFlags)

var configFile string
var config Config
var cacheFile string
var cache Cache

var firestoreClient *estimation_client.EstimationClient

var workflow alfred.Workflow

func main() {

    var err error
    if workflow, err = alfred.OpenWorkflow("..", true); err != nil {
        fmt.Printf("Error: %s", err)
        os.Exit(1)
    }

    dlog.Printf("cache file dir: %s", workflow.CacheDir())
    dlog.Printf("config file dir: %s", workflow.DataDir())

	configFile = path.Join(workflow.DataDir(), "config.json")
	if err := alfred.LoadJSON(configFile, &config); err != nil {
		dlog.Println("Error loading config:", err)
	}

    if config.APIKey == "" {
        fmt.Printf("APIKey is empty. Please write TOGGL_API_KEY to %s", configFile)
        os.Exit(1)
    }
    dlog.Printf("APIKey : %s", config.APIKey)

	cacheFile = path.Join(workflow.CacheDir(), "cache.json")
	if err := alfred.LoadJSON(cacheFile, &cache); err != nil {
		dlog.Println("Error loading cache:", err)
	}

    // firestore
    serviceAccount := option.WithCredentialsFile("../credential/secret.json")
    firestoreClient, err = estimation_client.Init(serviceAccount)
    if err != nil {
        log.Fatalln(err)
        os.Exit(1)
    }

    workflow.Run([]alfred.Command{
        AddEntryCommand{},
        GetEntryCommand{},
    })

    firestoreClient.Close()
}
