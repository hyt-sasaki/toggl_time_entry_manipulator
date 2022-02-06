package list

import (
	"fmt"
	"log"
	"os"
	"time"
    "strconv"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/config"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.list]", log.LstdFlags)

type ListEntryCommand struct {
    Repo repository.ICachedRepository
    Config config.WorkflowConfig
}

func NewListEntryCommand(repo repository.ICachedRepository, workflowConfig config.WorkflowConfig) (ListEntryCommand) {
    return ListEntryCommand{Repo: repo, Config: workflowConfig}
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
        projectAlias := config.GetAlias(c.Config.ProjectAliases, entity.Entry.Pid)
        // TODO tagのalias
        if !command.Match(title + projectAlias, arg) {
            continue
        }
        detailRefData := alfred.Stringify(command.DetailRefData{ID: entity.Entry.ID})
        item := alfred.Item{
            Title: getTitle(entity, projects),
            Subtitle: getSubtitle(entity),
            Icon: getIcon(entity),
            Arg: &alfred.ItemArg{
                Keyword: command.GetEntryKeyword,
                Mode: alfred.ModeTell,
                Data: detailRefData,
            },
        }
        item.AddMod(
            alfred.ModCmd,
            alfred.ItemMod{
                Subtitle: "Add this entry to favorite list",
                Arg: &alfred.ItemArg{
                    Keyword: command.FavoriteEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: detailRefData,
                },
            },
        )
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

// TODO testを追加
func getSubtitle(entity domain.TimeEntryEntity) (subtitle string) {
    if entity.IsRunning() {
        subtitle = fmt.Sprintf("%s ... (", convertTimeToString(entity.Entry.Start))
    } else {
        subtitle = fmt.Sprintf("%s - %s (actual: %s min, ", convertTimeToString(entity.Entry.Start), convertTimeToString(entity.Entry.Stop), convertDuration(entity.Entry.Duration))
    }
    if entity.HasEstimation() {
        subtitle = fmt.Sprintf("%sestimation: %d min)", subtitle, entity.Estimation.Duration)
    }

    return
}

func convertTimeToString(t *time.Time) (dateStr string) {
    if t == nil {
        return
    }
    ny, nm, nd := time.Now().In(time.Local).Date()
    ty, tm, td := t.In(time.Local).Date()
    layout := "06/01/02 15:04"
    if (ny == ty && nm == tm && nd  == td) {
        layout = "15:04"
    }
    dateStr = t.In(time.Local).Format(layout)
    return
}

func convertDuration(duration int64) string {
    if duration < 0 {
        return "-"
    }
    min := int(duration / 60)
    return strconv.Itoa(min)
}

func getIcon(entity domain.TimeEntryEntity) (icon string) {
    icon = command.OffIcon
    if entity.IsRunning() {
        icon = command.OnIcon
    }
    if entity.IsLate() {
        if entity.Estimation.Memo == "" {
            icon = command.LateIcon
        } else {
            icon = command.LateCheckedIcon
        }
    }
    return
}
