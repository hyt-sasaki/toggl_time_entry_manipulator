package delete


import (
    "fmt"
	"encoding/json"
	"log"
	"os"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/command"
	"github.com/jason0x43/go-alfred"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.delete]", log.LstdFlags)

type IDeleteEntryCommand interface {
    alfred.Action
}

type DeleteEntryCommand struct {
    repo repository.ICachedRepository
}

func NewDeleteEntryCommand(repo repository.ICachedRepository) (IDeleteEntryCommand) {
    return &DeleteEntryCommand{repo: repo}
}

func (c DeleteEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.DeleteEntryKeyword,
        Description: "delete entry",
        IsEnabled: true,
    }
}

func (c DeleteEntryCommand) Do(data string) (out string, err error) {
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

    err = c.repo.Delete(&entity)
    if err != nil {
        return
    }
    
    out = fmt.Sprintf("Entry has been deleted. Description: %s", entity.Entry.Description)
    return
}
