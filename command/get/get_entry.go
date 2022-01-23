package get


import (
	"encoding/json"
    "fmt"
	"log"
	"os"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/command"
	"github.com/jason0x43/go-alfred"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.get]", log.LstdFlags)


type GetEntryCommand struct {
    Repo repository.ICachedRepository
}

func (c GetEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.GetEntryKeyword,
        Description: "get entry",
        IsEnabled: true,
    }
}

func (c GetEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    var itemData command.DetailRefData

    err = json.Unmarshal([]byte(data), &itemData)
    if err != nil {
        dlog.Printf("Invalid data")
        return
    }

    entity, err := c.Repo.FindOneById(itemData.ID)
    if err != nil {
        item := alfred.Item{
            Title: "Something went wrong",
        }
        items = append(items, item)
        return
    }

    descriptionItem := alfred.Item{
        Title: fmt.Sprintf("Description: %s", entity.Entry.Description),
        Arg: &alfred.ItemArg{
            Keyword: command.GetEntryKeyword,   // TODO ModifyDescriptionKeywordを実装
            Mode: alfred.ModeTell,
        },
    }
    items = append(items, descriptionItem)
    stopItem := alfred.Item{
        Title: "Stop this entry",
        Arg: &alfred.ItemArg{
            Keyword: command.GetEntryKeyword,   // TODO StopKeyword
            Mode: alfred.ModeTell,      // TODO ModeDoに変える
        },
    }
    items = append(items, stopItem)
    return
}
