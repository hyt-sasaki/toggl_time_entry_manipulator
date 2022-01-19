package get


import (
	"log"
	"os"
	"toggl_time_entry_manipulator/repository"
	"github.com/jason0x43/go-alfred"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.get]", log.LstdFlags)


type GetEntryCommand struct {
    Repo repository.ICachedRepository
}

const GetEntryKeyword = "get_entry"

func (c GetEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: GetEntryKeyword,
        Description: "get entry",
        IsEnabled: true,
    }
}

func (c GetEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    // テキトウ
    item := alfred.Item{
        Title: "test",
        Subtitle: "test",
        Arg: &alfred.ItemArg{
            Keyword: GetEntryKeyword,
            Mode: alfred.ModeTell,
        },
    }
    items = append(items, item)
    return
}
