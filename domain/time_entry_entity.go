package domain


import (
    "strconv"
    "github.com/jason0x43/go-toggl"
    "toggl_time_entry_manipulator/estimation_client"
)

type TimeEntryEntity struct {
    Entry toggl.TimeEntry
    Estimation estimation_client.Estimation
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
        Estimation: estimation_client.Estimation{
            Duration: duration,
        },
    }
}

func (entity *TimeEntryEntity) UpdateTimeEntry(timeEntry toggl.TimeEntry) {
    entity.Entry = timeEntry
}
