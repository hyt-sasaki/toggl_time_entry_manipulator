package favorite

import (
    "fmt"
	"log"
	"os"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/config"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.list]", log.LstdFlags)

type FavoriteEntryCommand struct {
    Repo repository.ICachedRepository
    Config *config.WorkflowConfig
}

func NewFavoriteEntryCommand(repo repository.ICachedRepository, workflowConfig *config.WorkflowConfig) (FavoriteEntryCommand) {
    return FavoriteEntryCommand{Repo: repo, Config: workflowConfig}
}

func (c FavoriteEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.FavoriteEntryKeyword,
        Description: "favorite entries",
        IsEnabled: true,
    }
}

func (c FavoriteEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    projects, err := c.Repo.GetProjects()
    if err != nil {
        return
    }
    var mode alfred.ModeType
    if c.Config != nil && c.Config.RecordEstimate {
        mode = alfred.ModeTell
    } else {
        mode = alfred.ModeDo
    }
    for _, entityId := range c.Config.Favorites {
        entity, _ := c.Repo.FindOneById(entityId)
        title := getTitle(entity, projects)
        projectAlias := config.GetAlias(c.Config.ProjectAliases, entity.Entry.Pid)
        if !command.Match(title + projectAlias, arg) {
            continue
        }
        item := alfred.Item{
            Title: getTitle(entity, projects),
            Arg: &alfred.ItemArg{
                Keyword: command.ContinueEntryKeyword,
                Mode: mode,
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
