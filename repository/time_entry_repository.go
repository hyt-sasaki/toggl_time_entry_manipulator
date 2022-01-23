package repository

import (
	"log"
	"os"
    "toggl_time_entry_manipulator/domain"
    "toggl_time_entry_manipulator/client"
	"github.com/jason0x43/go-toggl"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.repository]", log.LstdFlags)

type Config struct {
	TogglAPIKey client.TogglApiKey `desc:"Toggl API key"`
}
type ConfigFile string


type ITimeEntryRepository interface {
    Fetch(toggl.Account) ([]domain.TimeEntryEntity, error)
    FetchTogglAccount() (toggl.Account, error)
    Insert(*domain.TimeEntryEntity) error
    Stop(*domain.TimeEntryEntity) error
}

type TimeEntryRepository struct {
    config *Config
    togglClient client.ITogglClient
    estimationClient client.IEstimationClient
}

func NewTimeEntryRepository(
    config *Config,
    togglClient client.ITogglClient,
    estimationClient client.IEstimationClient) (repo *TimeEntryRepository) {
    repo = &TimeEntryRepository{
        config: config,
        togglClient: togglClient,
        estimationClient: estimationClient,
    }
    return
}

func (repo *TimeEntryRepository) FetchTogglAccount() (account toggl.Account, err error) {
	account, err = repo.togglClient.GetAccount()
	if err != nil {
		return 
	}
    return
}

func (repo *TimeEntryRepository) Fetch(account toggl.Account) (entities []domain.TimeEntryEntity, err error) {
    entries := account.Data.TimeEntries
    var entryIds []int64
    for _, entry := range entries {
        entryIds = append(entryIds, int64(entry.ID))
    }
    estimations, err := repo.estimationClient.Fetch(entryIds)
    for idx, estimation := range estimations {
        entities = append(entities, domain.TimeEntryEntity{
            Entry: entries[idx],
            Estimation: estimation,
        })
    }
    return
}

func (repo *TimeEntryRepository) Insert(entity *domain.TimeEntryEntity) (err error) {
    entry, err := repo.togglClient.StartTimeEntry(entity.Entry.Description, entity.Entry.Pid, entity.Entry.Tags)
    entity.UpdateTimeEntry(entry)

    if err = repo.estimationClient.Insert(entity.GetId(), entity.Estimation); err != nil {
        return
    }

    return
}

func (repo *TimeEntryRepository) Stop(entity *domain.TimeEntryEntity) (err error) {
    entry, err := repo.togglClient.StopTimeEntry(entity.Entry)
    entity.UpdateTimeEntry(entry)

    return
}
