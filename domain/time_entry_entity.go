package domain


import (
    "time"
    "strconv"
    "github.com/jason0x43/go-toggl"
)

type TimeEntryEntity struct {
    Entry toggl.TimeEntry
    Estimation Estimation
}
type Estimation struct {
    Duration int        `firestore:"duration"`
    Memo string         `firestore:"memo"`
    CreatedTm time.Time `firestore:"createdTm"`
    UpdatedTm time.Time `firestore:"updatedTm"`
}

func (entity TimeEntryEntity) GetId() string {
    return strconv.Itoa(entity.Entry.ID)
}

func Create(description string, pid int, tag string, duration int) (entity TimeEntryEntity) {
    return TimeEntryEntity{
        Entry: toggl.TimeEntry{
            Pid: pid,
            Description: description,
            Tags: []string{tag},
        },
        Estimation: Estimation{
            Duration: duration,
            CreatedTm: time.Now(),
            UpdatedTm: time.Now(),
        },
    }
}

func (entity *TimeEntryEntity) UpdateTimeEntry(timeEntry toggl.TimeEntry) {
    entity.Entry = timeEntry
}
