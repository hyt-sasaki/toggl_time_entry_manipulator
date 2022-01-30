package stop

import (
    "fmt"
	"encoding/json"
	"log"
	"os"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/command"
	"github.com/jason0x43/go-alfred"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.get]", log.LstdFlags)

type StopEntryCommand struct {
    Repo repository.ICachedRepository
}

func NewStopEntryCommand(repo repository.ICachedRepository) (StopEntryCommand) {
    return StopEntryCommand{Repo: repo}
}

func (c StopEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.StopEntryKeyword,
        Description: "stop entry",
        IsEnabled: true,
    }
}

func (c StopEntryCommand) Do(data string) (out string, err error) {
    var itemData command.DetailRefData

    err = json.Unmarshal([]byte(data), &itemData)
    if err != nil {
        dlog.Printf("Invalid data")
        return
    }

    entity, err := c.Repo.FindOneById(itemData.ID)
    if err != nil {
        return
    }

    err = c.Repo.Stop(&entity)
    if err != nil {
        return
    }
    
    out = fmt.Sprintf("Entry has stopped. Description: %s", entity.Entry.Description)
    return
}
