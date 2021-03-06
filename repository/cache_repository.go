package repository

import (
    "fmt"
	"time"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository/myCache"

	"github.com/jason0x43/go-toggl"
)


type cachedRepository struct {
    cache myCache.ICache
    timeEntryRepository ITimeEntryRepository
}

type ICachedRepository interface {
    Fetch() ([]domain.TimeEntryEntity, error)
    FindOneById(int) (domain.TimeEntryEntity, error)
    GetProjects() ([]toggl.Project, error)
    GetTags() ([]toggl.Tag, error)
    Insert(*domain.TimeEntryEntity) (error)
    Update(*domain.TimeEntryEntity) (error)
    Stop(*domain.TimeEntryEntity) (error)
    Continue(*domain.TimeEntryEntity) (domain.TimeEntryEntity, error)
    Delete(*domain.TimeEntryEntity) (error)
}

func NewCachedRepository(
    cache myCache.ICache,
    timeEntryRepository ITimeEntryRepository) (repo ICachedRepository) {
    repo = &cachedRepository{
        cache: cache,
        timeEntryRepository: timeEntryRepository,
    }
    return
}

func (c *cachedRepository) Fetch() (entities []domain.TimeEntryEntity, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    entities = c.cache.GetData().Entities
    return
}

func (c *cachedRepository) FindOneById(entryId int) (entity domain.TimeEntryEntity, err error) {
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

func (c *cachedRepository) GetProjects() (projects []toggl.Project, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    projects = c.cache.GetData().Projects

    return
}

func (c *cachedRepository) GetTags() (tags []toggl.Tag, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    tags = c.cache.GetData().Tags

    return
}

func (c *cachedRepository) Insert(entity *domain.TimeEntryEntity) (err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    tags, _ := c.GetTags()
    if err = c.timeEntryRepository.Insert(entity, tags); err != nil {
        return
    }
	if err = c.refresh(); err != nil {
		return
	}
    return
}

func (c *cachedRepository) Update(entity *domain.TimeEntryEntity) (err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    tags, _ := c.GetTags()
    if err = c.timeEntryRepository.Update(entity, tags); err != nil {
        return
    }
	if err = c.refresh(); err != nil {
		return
	}
    return
}

func (c *cachedRepository) Stop(entity *domain.TimeEntryEntity) (err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    if err = c.timeEntryRepository.Stop(entity); err != nil {
        return
    }
	if err = c.refresh(); err != nil {
		return
	}
    return
}

func (c *cachedRepository) Delete(entity *domain.TimeEntryEntity) (err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    if err = c.timeEntryRepository.Delete(entity); err != nil {
        return
    }
	if err = c.refresh(); err != nil {
		return
	}
    return
}

func (c *cachedRepository) Continue(entity *domain.TimeEntryEntity) (newEntity domain.TimeEntryEntity, err error) {
	if err = c.checkRefresh(); err != nil {
		return
	}
    newEntity, err = c.timeEntryRepository.Continue(entity); 
    if err != nil {
        return
    }
	if err = c.refresh(); err != nil {
		return
	}
    return
}


func (c *cachedRepository) checkRefresh() error {
    t := c.cache.GetData().Time
	if time.Now().Sub(t).Minutes() < 1.0 {
		return nil
	}

	dlog.Println("Refreshing cache...")
	err := c.refresh()
	if err != nil {
		dlog.Println("Error refreshing cache:", err)
	}
	return err
}

func (c *cachedRepository) refresh() (err error) {
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
