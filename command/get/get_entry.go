package get

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"

	"github.com/jason0x43/go-alfred"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.add]", log.LstdFlags)

type GetEntryCommand struct {
    Repo repository.ICachedRepository
}

const GetEntryKeyword = "get_entry"

func (c GetEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: GetEntryKeyword,
        Description: "get entries",
        IsEnabled: true,
    }
}

func (c GetEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    entities, err := c.Repo.Fetch()
    for _, entity := range entities {
        if !filterByArg(arg, entity) {
            continue
        }
        item := alfred.Item{
            Title: fmt.Sprintf("Description: %s", entity.Entry.Description),
            Subtitle: fmt.Sprintf("actual duration: %s, estimation: %d", convertDuration(entity.Entry.Duration), entity.Estimation.Duration),
            Arg: &alfred.ItemArg{
                Keyword: GetEntryKeyword,
                Mode: alfred.ModeTell,
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
