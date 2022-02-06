package myCache

import (
	"time"
	"toggl_time_entry_manipulator/domain"

	"github.com/jason0x43/go-toggl"
)

type CacheFile string
type ICache interface {
    Save()
    GetData() *Data
}
type cache struct {
    data *Data
    file CacheFile
    saveCallback (func(string, interface{}) error)
}
type Data struct {
	Workspace int
    Projects  []toggl.Project
    Tags      []toggl.Tag
    Entities  []domain.TimeEntryEntity   
	Time      time.Time
}

func NewCache(data *Data, file CacheFile, callback func(string, interface{}) error) ICache {
    return &cache{
        data: data,
        file: file,
        saveCallback: callback,
    }
}

func (c *cache) Save() {
    c.saveCallback(string(c.file), c.data)
}

func (c *cache) GetData() *Data {
    return c.data
}
