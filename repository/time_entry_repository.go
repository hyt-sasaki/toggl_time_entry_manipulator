package repository

import (
    "time"
	"log"
	"os"
	"sort"
    "strconv"
	"toggl_time_entry_manipulator/client"
	"toggl_time_entry_manipulator/domain"

	"github.com/jason0x43/go-toggl"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.repository]", log.LstdFlags)


type ITimeEntryRepository interface {
    Fetch(toggl.Account) ([]domain.TimeEntryEntity, error)
    FetchTogglAccount() (toggl.Account, error)
    Insert(*domain.TimeEntryEntity) error
    Update(*domain.TimeEntryEntity) error
    Delete(*domain.TimeEntryEntity) error
    Stop(*domain.TimeEntryEntity) error
    Continue(*domain.TimeEntryEntity) (domain.TimeEntryEntity, error)
}

type TimeEntryRepository struct {
    togglClient client.ITogglClient
    estimationClient client.IEstimationClient
}

func NewTimeEntryRepository(
    togglClient client.ITogglClient,
    estimationClient client.IEstimationClient) (repo *TimeEntryRepository) {
    repo = &TimeEntryRepository{
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
    sort.SliceStable(entries, func(i, j int) bool {
        return entries[i].StartTime().After(entries[j].StartTime());
    })
    var entryIds []int64
    for _, entry := range entries {
        entryIds = append(entryIds, int64(entry.ID))
    }
    estimations, err := repo.estimationClient.Fetch(entryIds)
    for idx, estimation := range estimations {
        if estimation == nil {
            entities = append(entities, domain.TimeEntryEntity{
                Entry: entries[idx],
            })
        } else {
            entities = append(entities, domain.TimeEntryEntity{
                Entry: entries[idx],
                Estimation: *estimation,
            })
        }
    }
    return
}

func (repo *TimeEntryRepository) Insert(entity *domain.TimeEntryEntity) (err error) {
    entry, err := repo.togglClient.StartTimeEntry(entity.Entry.Description, entity.Entry.Pid, entity.Entry.Tags)
    entity.UpdateTimeEntry(entry)

    entity.Estimation.CreatedTm = time.Now()
    entity.Estimation.UpdatedTm = time.Now()
    if err = repo.estimationClient.Insert(entity.GetId(), entity.Estimation); err != nil {
        return
    }

    return
}

func (repo *TimeEntryRepository) Update(entity *domain.TimeEntryEntity) (err error) {
    entry, err := repo.togglClient.UpdateTimeEntry(entity.Entry)
    entity.UpdateTimeEntry(entry)

    if err = repo.estimationClient.Update(entity.GetId(), entity.Estimation); err != nil {
        return
    }

    return
}

func (repo *TimeEntryRepository) Stop(entity *domain.TimeEntryEntity) (err error) {
    entry, err := repo.togglClient.StopTimeEntry(entity.Entry)
    entity.UpdateTimeEntry(entry)

    return
}

func (repo *TimeEntryRepository) Continue(entity *domain.TimeEntryEntity) (newEntity domain.TimeEntryEntity, err error) {
    // TODO
    newEntry, err := repo.togglClient.ContinueTimeEntry(entity.Entry)
    if err != nil {
        return
    }
    id := strconv.Itoa(newEntry.ID)
    newEstimation := entity.Estimation.Copy()
    err = repo.estimationClient.Insert(id, newEstimation)   // TODO CreatedTm, UpdatedTmもnewEstimationに反映できるようにする
    if err != nil {
        repo.togglClient.DeleteTimeEntry(newEntry)
        return 
    }
    newEntity = domain.TimeEntryEntity{
        Entry: newEntry,
        Estimation: newEstimation,
    }
    return
}

func (repo *TimeEntryRepository) Delete(entity *domain.TimeEntryEntity) (err error) {
    err = repo.estimationClient.Delete(entity.GetId())
    if err != nil {
        return
    }
    err = repo.togglClient.DeleteTimeEntry(entity.Entry)
    if err != nil {
        rollbackFail := repo.estimationClient.Insert(entity.GetId(), entity.Estimation) // rollback
        if rollbackFail != nil {
            dlog.Printf("rollback failed for id = %d: %s", entity.Entry.ID, rollbackFail)
        }
        return
    }

    return
}
