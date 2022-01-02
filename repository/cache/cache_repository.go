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

type Cache struct {
	Workspace int
    Projects  []toggl.Project
    Tags      []toggl.Tag
    Entities  []domain.TimeEntryEntity   
	Time      time.Time
}
type CacheFile string

type CachedRepository struct {
    cache *Cache
    cacheFile CacheFile
    timeEntryRepository *repository.TimeEntryRepository
}

func NewCachedRepository(
    cache *Cache,
    cacheFile CacheFile,
    timeEntryRepository *repository.TimeEntryRepository) (repo *CachedRepository) {
    repo = &CachedRepository{
        cache: cache,
        cacheFile: cacheFile,
        timeEntryRepository: timeEntryRepository,
    }
    return
}

func (c *CachedRepository) Fetch() (entities []domain.TimeEntryEntity, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    entities = c.cache.Entities
    return
}

func (c *CachedRepository) GetProjects() (projects []toggl.Project, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    projects = c.cache.Projects

    return
}

func (c *CachedRepository) GetTags() (tags []toggl.Tag, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    tags = c.cache.Tags

    return
}

func (c *CachedRepository) checkRefresh() error {
	if time.Now().Sub(c.cache.Time).Minutes() < 5.0 {
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
		return err
	}
    entities, err := c.timeEntryRepository.Fetch(account)
    if err != nil {
        return err
    }

	dlog.Printf("got account: %#v", account)

	c.cache.Time = time.Now()
	c.cache.Projects = account.Data.Projects
    c.cache.Tags = account.Data.Tags
    c.cache.Entities = entities
	c.cache.Workspace = account.Data.Workspaces[0].ID

	return alfred.SaveJSON(string(c.cacheFile), &c.cache)
}
