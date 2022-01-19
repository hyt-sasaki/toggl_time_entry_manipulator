package repository

import (
    "fmt"
	"time"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository/myCache"

	"github.com/jason0x43/go-toggl"
)


type CachedRepository struct {
    cache myCache.ICache
    timeEntryRepository *TimeEntryRepository
}

type ICachedRepository interface {
    Fetch() ([]domain.TimeEntryEntity, error)
    FindOneById(int) (domain.TimeEntryEntity, error)
    GetProjects() ([]toggl.Project, error)
    GetTags() ([]toggl.Tag, error)
    Insert(*domain.TimeEntryEntity) (error)
}

func NewCachedRepository(
    cache myCache.ICache,
    timeEntryRepository *TimeEntryRepository) (repo *CachedRepository) {
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
    entities = c.cache.GetData().Entities
    return
}

// TODO test, mock追加
func (c *CachedRepository) FindOneById(entryId int) (entity domain.TimeEntryEntity, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    entities := c.cache.GetData().Entities
    
    for _, e := range entities {
        if e.Entry.ID == entryId {
            entity = e
            return
        }
    }

    err = fmt.Errorf("Resource not found: %d", entryId)

    return
}

func (c *CachedRepository) GetProjects() (projects []toggl.Project, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    projects = c.cache.GetData().Projects

    return
}

func (c *CachedRepository) GetTags() (tags []toggl.Tag, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    tags = c.cache.GetData().Tags

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
    t := c.cache.GetData().Time
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

    data := c.cache.GetData()
	data.Time = time.Now()
	data.Projects = account.Data.Projects
    data.Tags = account.Data.Tags
    data.Entities = entities
	data.Workspace = account.Data.Workspaces[0].ID
    c.cache.Save()

	return 
}
