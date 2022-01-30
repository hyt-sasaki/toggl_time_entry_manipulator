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

func (estimation *Estimation) Copy() (Estimation) {
    copied := *estimation
    return copied
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

// TODO test追加
func (entity *TimeEntryEntity) HasEstimation() bool {
    return !(entity.Estimation.CreatedTm.IsZero() && entity.Estimation.Duration == 0);
}

func (entity *TimeEntryEntity) IsRunning() bool {
    return entity.Entry.IsRunning();
}

func (entity *TimeEntryEntity) IsLate() bool {
    if (entity.IsRunning()) {
        return false
    }
    if (!entity.HasEstimation()) {
        return false
    }
    return entity.Entry.Duration > (int64)(entity.Estimation.Duration * 60)
}

func (entity *TimeEntryEntity) Copy() (TimeEntryEntity) {
    copied := TimeEntryEntity{}
    copied.Entry = entity.Entry.Copy()
    copied.Estimation = entity.Estimation.Copy()
    return copied
}
