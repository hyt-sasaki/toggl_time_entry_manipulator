package favorite

import (
    "fmt"
	"encoding/json"
	"log"
	"os"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/config"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.favorite]", log.LstdFlags)

type FavoriteEntryCommand struct {
    Repo repository.ICachedRepository
    Config *config.Config
    ConfigFile config.ConfigFile
}

func NewFavoriteEntryCommand(repo repository.ICachedRepository, config *config.Config, configFile config.ConfigFile) (FavoriteEntryCommand) {
    return FavoriteEntryCommand{Repo: repo, Config: config, ConfigFile: configFile}
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
    if c.Config != nil && c.Config.WorkflowConfig.RecordEstimate {
        mode = alfred.ModeTell
    } else {
        mode = alfred.ModeDo
    }
    for _, entityId := range c.Config.WorkflowConfig.Favorites {
        entity, findErr := c.Repo.FindOneById(entityId)
        if findErr != nil {
            continue
        }
        title := getTitle(entity, projects)
        projectAlias := config.GetAlias(c.Config.WorkflowConfig.ProjectAliases, entity.Entry.Pid)
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

func (c FavoriteEntryCommand) Do(data string) (out string, err error) {
    var itemData command.DetailRefData

    err = json.Unmarshal([]byte(data), &itemData)
    if err != nil {
        dlog.Printf("Invalid data")
        return
    }

    c.Config.WorkflowConfig.Favorites = append(c.Config.WorkflowConfig.Favorites, itemData.ID)
    alfred.SaveJSON(string(c.ConfigFile), *c.Config)
    out = "Entry has been added to favorite list."

    return
}


func getTitle(entity domain.TimeEntryEntity, projects []toggl.Project) (title string){
    projectName := "-"
    for _, p := range projects {
        if entity.Entry.Pid == p.ID {
            projectName = p.Name
        }
    }
    tags := entity.Entry.Tags
    if len(tags) == 0 {
        tags = append(tags, "No Tag")
    }
    title = fmt.Sprintf("%s (%s) %s", entity.Entry.Description, projectName, tags)
    return
}
