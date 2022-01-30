package command

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"toggl_time_entry_manipulator/config"
	"toggl_time_entry_manipulator/domain"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
	"golang.org/x/text/unicode/norm"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command]", log.LstdFlags)

func GenerateItemsForProject(
    projects []toggl.Project,
    arg string,
    entity domain.TimeEntryEntity,
    config config.WorkflowConfig,
    itemArgGenerator func(domain.TimeEntryEntity) alfred.ItemArg,
) (items []alfred.Item) {
    if arg == "" {
        for _, ac := range config.ProjectAutocompleteItems {
            item := alfred.Item{
                Title: ac,
                Subtitle: "For complete",
                Autocomplete: ac,
            }
            items = append(items, item)
        }
    }
    for _, project := range projects {
        if arg != "" {
            if !Match(project.Name, arg) {
                continue
            }
        }
        entity.Entry.Pid = project.ID
        itemArg := itemArgGenerator(entity)
        item := alfred.Item{
            Title: fmt.Sprintf("Project: %s", project.Name),
            Autocomplete: project.Name,
            Arg: &itemArg,
        }
        items = append(items, item)
    }
    return 
}

func GenerateItemsForTag(
    tags []toggl.Tag,
    arg string,
    entity domain.TimeEntryEntity,
    itemArgGenerator func(domain.TimeEntryEntity) alfred.ItemArg,
) (items []alfred.Item) {
    if arg == "" {
        itemArg := itemArgGenerator(entity)
        entity.Entry.Tags = []string{}
        noTagItem := alfred.Item{
            Title: "No tag",
            Arg: &itemArg,
        }
        items = append(items, noTagItem)
    }
    for _, tag := range tags {
        if arg != "" {
            if !Match(tag.Name, arg) {
                continue
            }
        }
        entity.Entry.Tags = []string{tag.Name}
        itemArg := itemArgGenerator(entity)
        item := alfred.Item{
            Title: fmt.Sprintf("Tag: %s", tag.Name),
            Autocomplete: tag.Name,
            Arg: &itemArg,
        }
        items = append(items, item)
    }
    return 
}

func GenerateItemsForEstimatedDuration(
    arg string,
    entity domain.TimeEntryEntity,
    itemArgGenerator func(domain.TimeEntryEntity) alfred.ItemArg,
) (items []alfred.Item) {
    autocomplete := fmt.Sprintf("%d", entity.Estimation.Duration)
    estimatedDuration, parseErr := strconv.Atoi(arg)
    icon := "power_off.png"
    var title string
    var itemArg *alfred.ItemArg
    if parseErr != nil {
        title = "Duration: -"
        itemArg = nil
    } else {
        entity.Estimation.Duration = estimatedDuration
        _itemArg := itemArgGenerator(entity)
        itemArg = &_itemArg
        icon = "power_on.png"
        title = fmt.Sprintf("Duration: %d", estimatedDuration)
    }
    items = append(items, alfred.Item{
        Title: title,
        Subtitle: "Enter estimated duration",
        Autocomplete: autocomplete,
        Icon: icon,
        Arg: itemArg,
    })

    return
}

func Match(target, query string) (bool) {
    normedQuery := norm.NFKC.String(query)
    normedTarget := norm.NFKC.String(target)
    for _, word := range strings.Split(normedQuery, " ") {
        if (!strings.Contains(normedTarget, word)) {
            return false
        }
    }
    return true
}

func GenerateBackItem(keyword, data string) (alfred.Item) {
    return alfred.Item{
        Title: "Back",
        Icon: "back.png",
        Arg: &alfred.ItemArg{
            Keyword: keyword,
            Mode: alfred.ModeTell,
            Data: data,
        },
    }
}
