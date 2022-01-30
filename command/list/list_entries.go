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
	"github.com/jason0x43/go-toggl"
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
    if err != nil {
        return
    }
    projects, err := c.Repo.GetProjects()
    if err != nil {
        return
    }

    for _, entity := range entities {
        title := getTitle(entity, projects)
        if !filterByArg(arg, title) {
            continue
        }
        item := alfred.Item{
            Title: getTitle(entity, projects),
            Subtitle: getSubtitle(entity),
            Icon: getIcon(entity),
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

func getTitle(entity domain.TimeEntryEntity, projects []toggl.Project) (title string){
    projectName := "-"
    for _, p := range projects {
        if entity.Entry.Pid == p.ID {
            projectName = p.Name
        }
    }
    title = fmt.Sprintf("%s (%s)", entity.Entry.Description, projectName)
    return
}

func getSubtitle(entity domain.TimeEntryEntity) (subtitle string) {
    if entity.HasEstimation() {
        subtitle = fmt.Sprintf("actual duration: %s [min], estimation: %d [min]", convertDuration(entity.Entry.Duration), entity.Estimation.Duration)
    } else {
        subtitle = fmt.Sprintf("actual duration: %s [min], estimation: -", convertDuration(entity.Entry.Duration))
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

func filterByArg(arg, title string) (res bool) {
    res = true
    if arg == "" {
        return
    }
    args := strings.Split(arg, " ")
    for _, a := range(args) {
        res = res && strings.Contains(title, a)
    }
    return
}

func getIcon(entity domain.TimeEntryEntity) (icon string) {
    icon = "power_off.png"
    if entity.IsRunning() {
        icon = "power_on.png"
    }
    if entity.IsLate() {
        if entity.Estimation.Memo == "" {
            icon = "late.png"
        } else {
            icon = "late_checked.png"
        }
    }
    return
}
