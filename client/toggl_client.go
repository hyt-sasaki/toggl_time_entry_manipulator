package client

import (
	"log"
	"os"
	"toggl_time_entry_manipulator/config"

	"github.com/jason0x43/go-toggl"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.client]", log.LstdFlags)

type togglClient struct {
    config config.TogglConfig
}
type ITogglClient interface {
    GetAccount() (toggl.Account, error)
    StartTimeEntry(string, int, []string) (toggl.TimeEntry, error)
    StopTimeEntry(toggl.TimeEntry) (toggl.TimeEntry, error)
    ContinueTimeEntry(toggl.TimeEntry) (toggl.TimeEntry, error)
    DeleteTimeEntry(toggl.TimeEntry) (error)
    UpdateTimeEntry(toggl.TimeEntry) (toggl.TimeEntry, error)
}

func NewTogglClient(config config.TogglConfig) (ITogglClient) {
    return &togglClient{
        config: config,
    }
}

func (c *togglClient) GetAccount() (account toggl.Account, err error) {
	s := c.getSession()
	account, err = s.GetAccount()
	if err != nil {
		return 
	}
    return
}

func (c *togglClient) StartTimeEntry(description string, pid int, tags []string) (entry toggl.TimeEntry, err error) {
    s := c.getSession()

    entry, err = s.StartTimeEntryForProject(description, pid, false)
    if err != nil {
        return
    }

    if len(tags) > 0 {
        entry.Tags = tags
        _, err = s.UpdateTimeEntry(entry)
        if err != nil {
            return
        }
    }

    return
}

func (c *togglClient) StopTimeEntry(entry toggl.TimeEntry) (resultEntry toggl.TimeEntry, err error) {
    s := c.getSession()

    resultEntry, err = s.StopTimeEntry(entry)

    return
}

func (c *togglClient) ContinueTimeEntry(entry toggl.TimeEntry) (resultEntry toggl.TimeEntry, err error) {
    s := c.getSession()

    resultEntry, err = s.ContinueTimeEntry(entry, false)

    return
}

func (c *togglClient) UpdateTimeEntry(entry toggl.TimeEntry) (resultEntry toggl.TimeEntry, err error) {
    s := c.getSession()

    resultEntry, err = s.UpdateTimeEntry(entry)

    return
}

func (c *togglClient) DeleteTimeEntry(entry toggl.TimeEntry) (err error) {
    s := c.getSession()

    _, err = s.DeleteTimeEntry(entry)

    return
}

func (c *togglClient) getSession() (toggl.Session) {
    return toggl.OpenSession(string(c.config.APIKey))
}
