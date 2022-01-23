package list

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"
    "toggl_time_entry_manipulator/command"

	"github.com/jason0x43/go-alfred"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.list]", log.LstdFlags)

type ListEntryCommand struct {
    Repo repository.ICachedRepository
}


func (c ListEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.ListEntryKeyword,
        Description: "get entries",
        IsEnabled: true,
    }
}

func (c ListEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    entities, err := c.Repo.Fetch()
    for _, entity := range entities {
        if !filterByArg(arg, entity) {
            continue
        }
        item := alfred.Item{
            Title: fmt.Sprintf("Description: %s", entity.Entry.Description),
            Subtitle: fmt.Sprintf("actual duration: %s [min], estimation: %d [min]", convertDuration(entity.Entry.Duration), entity.Estimation.Duration),
            Arg: &alfred.ItemArg{
                Keyword: command.GetEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(command.DetailRefData{ID: entity.Entry.ID}),
            },
        }
        items = append(items, item)
    }
    return
}

func convertDuration(duration int64) string {
    if duration < 0 {
        return "[stil running...]"
    }
    min := int(duration / 60)
    return strconv.Itoa(min)
}

func filterByArg(arg string, entity domain.TimeEntryEntity) (res bool) {
    if arg == "" {
        res = true
        return
    }
    res = strings.Contains(entity.Entry.Description, arg)
    return
}
