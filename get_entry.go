package main

import (
    "fmt"
    "strconv"
    "github.com/jason0x43/go-alfred"
    cacheRepo "toggl_time_entry_manipulator/repository/cache"
)

type GetEntryCommand struct {
    repo *cacheRepo.CachedRepository
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
    dlog.Printf("Items")
    entities, err := c.repo.Fetch()
    for _, entity := range entities {
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
