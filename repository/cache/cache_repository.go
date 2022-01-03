package cache

import (
	"time"
    "os"
    "log"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.cache]", log.LstdFlags)

type CacheFile string
type Cache struct {
    Data *Data
    File CacheFile
}
type Data struct {
	Workspace int
    Projects  []toggl.Project
    Tags      []toggl.Tag
    Entities  []domain.TimeEntryEntity   
	Time      time.Time
}

func (c *Cache) Save() {
	alfred.SaveJSON(string(c.File), &c.Data)
}

type CachedRepository struct {
    cache *Cache
    timeEntryRepository *repository.TimeEntryRepository
}

func NewCachedRepository(
    cache *Cache,
    timeEntryRepository *repository.TimeEntryRepository) (repo *CachedRepository) {
    repo = &CachedRepository{
        cache: cache,
        timeEntryRepository: timeEntryRepository,
    }
    return
}

func (c *CachedRepository) Fetch() (entities []domain.TimeEntryEntity, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    entities = c.cache.Data.Entities
    return
}

func (c *CachedRepository) GetProjects() (projects []toggl.Project, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    projects = c.cache.Data.Projects

    return
}

func (c *CachedRepository) GetTags() (tags []toggl.Tag, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    tags = c.cache.Data.Tags

    return
}

func (c *CachedRepository) Insert(entity *domain.TimeEntryEntity) (err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    if err = c.timeEntryRepository.Insert(entity); err != nil {
        return
    }
	if err = c.checkRefresh(); err != nil {
		return
	}
    return
}


func (c *CachedRepository) checkRefresh() error {
    t := c.cache.Data.Time
	if time.Now().Sub(t).Minutes() < 5.0 {
		return nil
	}

	dlog.Println("Refreshing cache...")
	err := c.refresh()
	if err != nil {
		dlog.Println("Error refreshing cache:", err)
	}
	return err
}

func (c *CachedRepository) refresh() (err error) {
	account, err := c.timeEntryRepository.FetchTogglAccount()
	if err != nil {
		return
	}
    entities, err := c.timeEntryRepository.Fetch(account)
    if err != nil {
        return
    }

	dlog.Printf("got account: %#v", account)

	c.cache.Data.Time = time.Now()
	c.cache.Data.Projects = account.Data.Projects
    c.cache.Data.Tags = account.Data.Tags
    c.cache.Data.Entities = entities
	c.cache.Data.Workspace = account.Data.Workspaces[0].ID
    c.cache.Save()

	return 
}
