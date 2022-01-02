package estimation_client

import (
	"github.com/jason0x43/go-toggl"
)

type TogglApiKey string
type TogglClient struct {
    apiKey TogglApiKey
}
type ITogglClient interface {
    GetAccount() (toggl.Account, error)
    StartTimeEntry(string, int, []string) error
}

func NewTogglClient(apiKey TogglApiKey) (*TogglClient) {
    return &TogglClient{
        apiKey: apiKey,
    }
}

func (c *TogglClient) GetAccount() (account toggl.Account, err error) {
	s := c.getSession()
	account, err = s.GetAccount()
	if err != nil {
		return 
	}
    return
}

func (c *TogglClient) StartTimeEntry(description string, pid int, tags []string) (err error) {
    s := c.getSession()

    entry, err := s.StartTimeEntryForProject(description, pid, false)
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

func (c *TogglClient) getSession() (toggl.Session) {
    return toggl.OpenSession(string(c.apiKey))
}
