package main

import (
    "fmt"
    "os"
	"path"
    "errors"
	"google.golang.org/api/option"
	"github.com/jason0x43/go-alfred"
    "toggl_time_entry_manipulator/config"
    "toggl_time_entry_manipulator/repository/myCache"
    "toggl_time_entry_manipulator/command/add"
    "toggl_time_entry_manipulator/command/list"
    "toggl_time_entry_manipulator/command/get"
    "toggl_time_entry_manipulator/command/modify"
    "toggl_time_entry_manipulator/command/stop"
    "toggl_time_entry_manipulator/command/delete"
    "toggl_time_entry_manipulator/command/continue_entry"
)

const configFileName = "config.json"
const cacheFileName = "cache.json"

func NewServiceAccount(workflow alfred.Workflow) (serviceAccount option.ClientOption, err error) {
    filePath := path.Join(workflow.DataDir(), "secret.json")
    if !exists(filePath) {
        err = fmt.Errorf("%s does not exist.", filePath)
        return
    }
    serviceAccount = option.WithCredentialsFile(filePath)
    return
}

func NewConfigFile(workflow alfred.Workflow) config.ConfigFile {
    configFile := path.Join(workflow.DataDir(), configFileName)
    return config.ConfigFile(configFile)
}

func NewConfig(configFile config.ConfigFile) (config *config.Config, err error) {
	if err = alfred.LoadJSON(string(configFile), &config); err != nil {
		dlog.Println("Error loading config:", err)
        return
	}
    if config.TogglConfig.APIKey == "" {
        dlog.Printf("APIKey is empty. Please write TogglConfig.APIKey to %s", configFile)
        err = fmt.Errorf("APIKey is empty. Please write TogglConfig.APIKey to %s", configFile)
        return
    }
    if config.FirestoreConfig.CollectionName == "" {
        dlog.Printf("Firestore collection name is empty. Please write Firestore.CollectionName to %s", configFile)
        err = fmt.Errorf("CollectionName is empty. Please write Firestore.CollectionName to %s", configFile)
        return
    }

    return 
}

func NewCacheFile(workflow alfred.Workflow) myCache.CacheFile {
    cacheFile := path.Join(workflow.CacheDir(), cacheFileName)
    return myCache.CacheFile(cacheFile)
}

func NewCache(cacheFile myCache.CacheFile) (cache *myCache.Cache, err error) {
    var data *myCache.Data
	if err = alfred.LoadJSON(string(cacheFile), &data); err != nil {
		dlog.Println("No cache file found:", err)
        data = &myCache.Data{}
        alfred.SaveJSON(string(cacheFile), *data)
	}
    cache = &myCache.Cache{
        Data: data,
        File: cacheFile,
        SaveCallback: alfred.SaveJSON,
    }
    dlog.Println(cache)

    return
}

func NewCommands(
    firstCall bool,
    addCommand add.AddEntryCommand,
    listCommand list.ListEntryCommand,
    getCommand get.GetEntryCommand,
    modifyComamnd modify.ModifyEntryCommand,
    stopCommand stop.StopEntryCommand,
    deleteCommand delete.DeleteEntryCommand,
    continueCommand continue_entry.ContinueEntryCommand,
) []alfred.Command {
    if firstCall {
        return []alfred.Command{
            addCommand,
            listCommand,
        }
    } else {
        return []alfred.Command{
            addCommand,
            listCommand,
            getCommand,
            modifyComamnd,
            stopCommand,
            deleteCommand,
            continueCommand,
        }
    }
}


func exists(path string) bool {
    _, err := os.Stat(path)
    return !errors.Is(err, os.ErrNotExist)
}
