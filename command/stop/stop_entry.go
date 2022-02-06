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

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.stop]", log.LstdFlags)

type IStopEntryCommand interface {
    alfred.Action
}

type StopEntryCommand struct {
    repo repository.ICachedRepository
}

func NewStopEntryCommand(repo repository.ICachedRepository) (com IStopEntryCommand) {
    com = &StopEntryCommand{repo: repo}
    return
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

    entity, err := c.repo.FindOneById(itemData.ID)
    if err != nil {
        return
    }

    err = c.repo.Stop(&entity)
    if err != nil {
        return
    }
    
    out = fmt.Sprintf("Entry has stopped. Description: %s", entity.Entry.Description)
    return
}
