package main

import (
	"path"
	"github.com/jason0x43/go-alfred"
    "toggl_time_entry_manipulator/repository"
    cacheRepo "toggl_time_entry_manipulator/repository/cache"
)

const configFileName = "config.json"
const cacheFileName = "cache.json"

func NewConfigFile(workflow alfred.Workflow) repository.ConfigFile {
    configFile := path.Join(workflow.DataDir(), configFileName)
    return repository.ConfigFile(configFile)
}

func NewConfig(configFile repository.ConfigFile) (config *repository.Config, err error) {
	if err := alfred.LoadJSON(string(configFile), &config); err != nil {
		dlog.Println("Error loading config:", err)
	}

    return 
}

func NewCacheFile(workflow alfred.Workflow) cacheRepo.CacheFile {
    cacheFile := path.Join(workflow.CacheDir(), cacheFileName)
    return cacheRepo.CacheFile(cacheFile)
}

func NewCache(cacheFile cacheRepo.CacheFile) (cache *cacheRepo.Cache, err error) {
	if err = alfred.LoadJSON(string(cacheFile), &cache); err != nil {
		dlog.Println("Error loading cache:", err)
        return
	}

    return
}
