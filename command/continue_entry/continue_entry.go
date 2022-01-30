package continue_entry

import (
    "fmt"
	"encoding/json"
	"log"
	"os"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/command"
	"github.com/jason0x43/go-alfred"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.continue]", log.LstdFlags)

type ContinueEntryCommand struct {
    Repo repository.ICachedRepository
}

func (c ContinueEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.ContinueEntryKeyword,
        Description: "continue entry",
        IsEnabled: true,
    }
}

func (c ContinueEntryCommand) Do(data string) (out string, err error) {
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

    newEntity, err := c.Repo.Continue(&entity)
    if err != nil {
        return
    }
    
    out = fmt.Sprintf("Entry has been copied. Description: %s", newEntity.Entry.Description)
    return
}
