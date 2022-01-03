package myCache

import (
	"time"
	"toggl_time_entry_manipulator/domain"

	//"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
)

type CacheFile string
type ICache interface {
    Save()
    GetData() *Data
}
type Cache struct {
    Data *Data
    File CacheFile
    SaveCallback (func(string, interface{}) error)
}
type Data struct {
	Workspace int
    Projects  []toggl.Project
    Tags      []toggl.Tag
    Entities  []domain.TimeEntryEntity   
	Time      time.Time
}

func (c *Cache) Save() {
    c.SaveCallback(string(c.File), c.Data)
	// alfred.SaveJSON(string(c.File), &c.Data)
}

func (c *Cache) GetData() *Data {
    return c.Data
}
