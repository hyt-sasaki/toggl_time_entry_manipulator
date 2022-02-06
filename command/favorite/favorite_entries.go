package favorite

import (
    "fmt"
	"encoding/json"
	"log"
	"os"
    "errors"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/config"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.favorite]", log.LstdFlags)

type IFavoriteEntryCommand interface {
    alfred.Filter
    alfred.Action
}

type FavoriteEntryCommand struct {
    repo repository.ICachedRepository
    config *config.Config
    configFile config.ConfigFile
}

func NewFavoriteEntryCommand(repo repository.ICachedRepository, config *config.Config, configFile config.ConfigFile) (com IFavoriteEntryCommand, err error) {
    if config == nil {
        err = errors.New("Config is nil.")
        return
    }

    com = &FavoriteEntryCommand{repo: repo, config: config, configFile: configFile}
    return 
}

func (c FavoriteEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.FavoriteEntryKeyword,
        Description: "favorite entries",
        IsEnabled: true,
    }
}

func (c FavoriteEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    projects, err := c.repo.GetProjects()
    if err != nil {
        return
    }
    var mode alfred.ModeType
    if c.config != nil && c.config.WorkflowConfig.RecordEstimate {
        mode = alfred.ModeTell
    } else {
        mode = alfred.ModeDo
    }
    for _, entityId := range c.config.WorkflowConfig.Favorites {
        entity, findErr := c.repo.FindOneById(entityId)
        if findErr != nil {
            continue
        }
        title := getTitle(entity, projects)
        projectAlias := config.GetAlias(c.config.WorkflowConfig.ProjectAliases, entity.Entry.Pid)
        if !command.Match(title + projectAlias, arg) {
            continue
        }
        detailRefData := command.DetailRefData{ID: entity.Entry.ID}
        item := alfred.Item{
            Title: getTitle(entity, projects),
            Arg: &alfred.ItemArg{
                Keyword: command.ContinueEntryKeyword,
                Mode: mode,
                Data: alfred.Stringify(detailRefData),
            },
        }
        item.AddMod(alfred.ModCmd, alfred.ItemMod{
            Subtitle: "Remove this entry from favorite list",
            Arg: &alfred.ItemArg{
                Keyword: command.FavoriteEntryKeyword,
                Mode: alfred.ModeDo,
                Data: alfred.Stringify(command.FavoriteRefData{
                    Ref: detailRefData,
                    Action: command.RemoveFromFavorite,
                }),
            },
        })
        items = append(items, item)
    }

    return
}

func (c FavoriteEntryCommand) Do(data string) (out string, err error) {
    var itemData command.FavoriteRefData

    err = json.Unmarshal([]byte(data), &itemData)
    if err != nil {
        dlog.Printf("Invalid data")
        return
    }

    switch itemData.Action {
    case command.AddToFavorite:
        c.config.WorkflowConfig.Favorites = append(c.config.WorkflowConfig.Favorites, itemData.Ref.ID)
        alfred.SaveJSON(string(c.configFile), *c.config)
        out = "Entry has been added to favorite list."
    case command.RemoveFromFavorite:
        removeIdx := -1
        for idx, entryId := range c.config.WorkflowConfig.Favorites {
            if entryId == itemData.Ref.ID {
                removeIdx = idx
            }
        }
        if removeIdx != -1 {
            c.config.WorkflowConfig.Favorites = append(c.config.WorkflowConfig.Favorites[:removeIdx], c.config.WorkflowConfig.Favorites[removeIdx+1:]...)
            alfred.SaveJSON(string(c.configFile), *c.config)
            out = "Entry has been removed from favorite list."
        } else {
            out = fmt.Sprintf("Not found %d in favorite list.", itemData.Ref.ID)
        } 
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
    tags := entity.Entry.Tags
    if len(tags) == 0 {
        tags = append(tags, "No Tag")
    }
    title = fmt.Sprintf("%s (%s) %s", entity.Entry.Description, projectName, tags)
    return
}
