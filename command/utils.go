package command

import (
    "fmt"
    "strings"
	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
    "toggl_time_entry_manipulator/domain"
)

func GenerateItemsForProject(
    projects []toggl.Project,
    arg string,
    entity domain.TimeEntryEntity,
    itemArgGenerator func(domain.TimeEntryEntity) alfred.ItemArg,
) (items []alfred.Item) {
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

func Match(target, query string) (bool) {
    return strings.Contains(target, query)
}
