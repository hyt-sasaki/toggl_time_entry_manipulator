package main

import (
    "fmt"
	"path"
	"github.com/jason0x43/go-alfred"
    "toggl_time_entry_manipulator/repository"
    "toggl_time_entry_manipulator/repository/myCache"
)

const configFileName = "config.json"
const cacheFileName = "cache.json"

func NewConfigFile(workflow alfred.Workflow) repository.ConfigFile {
    configFile := path.Join(workflow.DataDir(), configFileName)
    return repository.ConfigFile(configFile)
}

func NewConfig(configFile repository.ConfigFile) (config *repository.Config, err error) {
	if err = alfred.LoadJSON(string(configFile), &config); err != nil {
		dlog.Println("Error loading config:", err)
        return
	}
    if config.TogglAPIKey == "" {
        dlog.Printf("APIKey is empty. Please write TOGGL_API_KEY to %s", configFile)
        err = fmt.Errorf("APIKey is empty. Please write TOGGL_API_KEY to %s", configFile)
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
		dlog.Println("Error loading cache:", err)
        return
	}
    cache = &myCache.Cache{
        Data: data,
        File: cacheFile,
        SaveCallback: alfred.SaveJSON,
    }
    dlog.Println(cache)

    return
}
