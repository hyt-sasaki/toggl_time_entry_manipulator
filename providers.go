package main

import (
    "os"
	"path"
    "errors"
	"google.golang.org/api/option"
	"github.com/jason0x43/go-alfred"
    workflowConfig "toggl_time_entry_manipulator/config"
    "toggl_time_entry_manipulator/repository/myCache"
    "toggl_time_entry_manipulator/command/add"
    "toggl_time_entry_manipulator/command/list"
    "toggl_time_entry_manipulator/command/favorite"
    "toggl_time_entry_manipulator/command/get"
    "toggl_time_entry_manipulator/command/modify"
    "toggl_time_entry_manipulator/command/stop"
    "toggl_time_entry_manipulator/command/delete"
    "toggl_time_entry_manipulator/command/continue_entry"
    optionCom "toggl_time_entry_manipulator/command/option"
)

const configFileName = "config.json"
const cacheFileName = "cache.json"

func NewServiceAccount(workflow alfred.Workflow) (serviceAccount option.ClientOption, err error) {
    filePath := path.Join(workflow.DataDir(), "secret.json")
    if !exists(filePath) {
		dlog.Printf("%s does not exist.\n", filePath)
        return
    }
    serviceAccount = option.WithCredentialsFile(filePath)
    return
}

func NewConfigFile(workflow alfred.Workflow) workflowConfig.ConfigFile {
    configFile := path.Join(workflow.DataDir(), configFileName)
    return workflowConfig.ConfigFile(configFile)
}

func NewConfig(configFile workflowConfig.ConfigFile) (config *workflowConfig.Config, err error) {
	if err = alfred.LoadJSON(string(configFile), &config); err != nil {
		dlog.Println("No cache file found:", err)
        config = &workflowConfig.Config{}
        config.WorkflowConfig.RecordEstimate = false
        alfred.SaveJSON(string(configFile), *config)
	}
    if config.TogglConfig.APIKey == "" {
        dlog.Printf("APIKey is empty. Please write TogglConfig.APIKey to %s", configFile)
    }
    if config.FirestoreConfig.CollectionName == "" && config.WorkflowConfig.RecordEstimate {
        dlog.Printf("Firestore collection name is empty. Please write Firestore.CollectionName to %s", configFile)
        return
    }

    return 
}

func NewCacheFile(workflow alfred.Workflow) myCache.CacheFile {
    cacheFile := path.Join(workflow.CacheDir(), cacheFileName)
    return myCache.CacheFile(cacheFile)
}

func NewCache(cacheFile myCache.CacheFile) (cache myCache.ICache, err error) {
    var data *myCache.Data
	if err = alfred.LoadJSON(string(cacheFile), &data); err != nil {
		dlog.Println("No cache file found:", err)
        data = &myCache.Data{}
        alfred.SaveJSON(string(cacheFile), *data)
	}
    cache = myCache.NewCache(data, cacheFile, alfred.SaveJSON)
    dlog.Println(cache)

    return
}

func NewCommands(
    firstCall bool,
    config *workflowConfig.Config,
    optionCommand optionCom.OptionCommand,
    addCommand add.IAddEntryCommand,
    listCommand list.IListEntryCommand,
    favoriteCommand favorite.FavoriteEntryCommand,
    getCommand get.IGetEntryCommand,
    modifyComamnd modify.ModifyEntryCommand,
    stopCommand stop.StopEntryCommand,
    deleteCommand delete.DeleteEntryCommand,
    continueCommand continue_entry.ContinueEntryCommand,
) []alfred.Command {
    if config.TogglConfig.APIKey == "" {
        return []alfred.Command{optionCommand}
    }
    if firstCall {
        return []alfred.Command{
            addCommand,
            listCommand,
            favoriteCommand,
            optionCommand,
        }
    } else {
        return []alfred.Command{
            addCommand,
            listCommand,
            favoriteCommand,
            getCommand,
            modifyComamnd,
            stopCommand,
            deleteCommand,
            continueCommand,
            optionCommand,
        }
    }
}


func exists(path string) bool {
    _, err := os.Stat(path)
    return !errors.Is(err, os.ErrNotExist)
}
