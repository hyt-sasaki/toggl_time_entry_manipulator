package get

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/repository"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.get]", log.LstdFlags)


type GetEntryCommand struct {
    Repo repository.ICachedRepository
}

func (c GetEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.GetEntryKeyword,
        Description: "get entry",
        IsEnabled: true,
    }
}

func (c GetEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    var itemData command.DetailRefData

    err = json.Unmarshal([]byte(data), &itemData)
    if err != nil {
        dlog.Printf("Invalid data")
        return
    }

    entity, err := c.Repo.FindOneById(itemData.ID)
    if err != nil {
        item := alfred.Item{
            Title: "Something went wrong",
        }
        items = append(items, item)
        return
    }
    projects, err := c.Repo.GetProjects()
    if err != nil {
        item := alfred.Item{
            Title: "Something went wrong",
        }
        items = append(items, item)
        return
    }

    if alfred.FuzzyMatches("description", arg) {
        descriptionItem := alfred.Item{
            Title: fmt.Sprintf("Description: %s", entity.Entry.Description),
            Arg: &alfred.ItemArg{
                Keyword: command.ModifyEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(command.ModifyData{
                    Ref: command.DetailRefData{ID: entity.Entry.ID},
                    Target: command.ModifyDescription,
                }),
            },
        }
        items = append(items, descriptionItem)
    }

    if alfred.FuzzyMatches("project", arg) {
        projectItem := alfred.Item{
            Title: fmt.Sprintf("Project: %s", getProject(entity.Entry.Pid, projects).Name),
            Arg: &alfred.ItemArg{
                Keyword: command.ModifyEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(command.ModifyData{
                    Ref: command.DetailRefData{ID: entity.Entry.ID},
                    Target: command.ModifyProject,
                }),
            },
        }
        items = append(items, projectItem)
    }

    if alfred.FuzzyMatches("tag", arg) {
        tagItem := alfred.Item{
            Title: fmt.Sprintf("Tag: %s", entity.Entry.Tags),
            Arg: &alfred.ItemArg{
                Keyword: command.ModifyEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(command.ModifyData{
                    Ref: command.DetailRefData{ID: entity.Entry.ID},
                    Target: command.ModifyTag,
                }),
            },
        }
        items = append(items, tagItem)
    }

    if entity.HasEstimation() && alfred.FuzzyMatches("estimated duration", arg) {
        estimatedDurationItem := alfred.Item{
            Title: fmt.Sprintf("Estimated duration: %d [min]", entity.Estimation.Duration),
            Arg: &alfred.ItemArg{
                Keyword: command.ModifyEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(command.ModifyData{
                    Ref: command.DetailRefData{ID: entity.Entry.ID},
                    Target: command.ModifyDuration,
                }),
            },
        }
        items = append(items, estimatedDurationItem)
    }

    timeLayout := "06/01/02 15:04"
    if alfred.FuzzyMatches("start", arg) {
        startTimeItem := alfred.Item{
            Title: fmt.Sprintf("Start: %s", entity.Entry.Start.In(time.Local).Format(timeLayout)),
            Arg: &alfred.ItemArg{
                Keyword: command.ModifyEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(command.ModifyData{
                    Ref: command.DetailRefData{ID: entity.Entry.ID},
                    Target: command.ModifyStart,
                }),
            },
        }
        items = append(items, startTimeItem)
    }
    if entity.Entry.Stop != nil && alfred.FuzzyMatches("stop", arg) {
        stopTimeItem := alfred.Item{
            Title: fmt.Sprintf("Stop: %s", entity.Entry.Stop.In(time.Local).Format(timeLayout)),
            Arg: &alfred.ItemArg{
                Keyword: command.ModifyEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(command.ModifyData{
                    Ref: command.DetailRefData{ID: entity.Entry.ID},
                    Target: command.ModifyStop,
                }),
            },
        }
        items = append(items, stopTimeItem)
    }

    if entity.HasEstimation() && alfred.FuzzyMatches("memo", arg) {
        memoItem := alfred.Item{
            Title: fmt.Sprintf("Memo: %s", entity.Estimation.Memo),
            Arg: &alfred.ItemArg{
                Keyword: command.ModifyEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(command.ModifyData{
                    Ref: command.DetailRefData{ID: entity.Entry.ID},
                    Target: command.ModifyMemo,
                }),
            },
        }
        items = append(items, memoItem)
    }

    if entity.IsRunning() && alfred.FuzzyMatches("stop this entry", arg) {
        stopItem := alfred.Item{
            Title: "Stop this entry",
            Arg: &alfred.ItemArg{
                Keyword: command.StopEntryKeyword,
                Mode: alfred.ModeDo,
                Data: data,
            },
        }
        items = append(items, stopItem)
    }

    return
}

func getProject(pid int, projects []toggl.Project) (project toggl.Project) {
    for _, p := range projects {
        if p.ID == pid {
            project = p
            return
        }
    }
    return
}
