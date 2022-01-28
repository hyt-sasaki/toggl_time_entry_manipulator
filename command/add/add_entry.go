package add

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.add]", log.LstdFlags)

type AddEntryCommand struct {
    Repo repository.ICachedRepository
}

type StateData struct {
    Current  state
    Args EntryArgs
}
type state int
const (
    DescriptionEdit state = iota
    ProjectEdit
    TagEdit 
    TimeEstimationEdit
    EndEdit
)
const initialState = ProjectEdit
type EntryArgs struct {
    Description string
    Project int
    Tag string
    TimeEstimation int  // minutes
}

func (c AddEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.AddEntryKeyword,
        Description: "add entry: project -> tag -> description -> estimation",
        IsEnabled: true,
    }
}

func (c AddEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    // load from alfred variable
    var sd StateData

	if data != "" {
		if err := json.Unmarshal([]byte(data), &sd); err != nil {
			dlog.Printf("Invalid data")
		}
	} else {
        sd = StateData{
            Current: initialState,
        }
    }


    // generate items
    switch sd.Current {
        case DescriptionEdit:
            items = append(items, c.generateDescriptionItems(sd, arg)...)
        case ProjectEdit:
            var projects []toggl.Project
            projects, err = c.Repo.GetProjects()
            if err != nil {
                return
            }
            items = append(items, c.generateProjectItems(sd, arg, projects)...)
        case TagEdit:
            var tags []toggl.Tag
            tags, err = c.Repo.GetTags()
            if err != nil {
                return
            }
            items = append(items, c.generateTagItems(sd, arg, tags)...)
        case TimeEstimationEdit:
            items = append(items, c.generateTimeEstimationItems(sd, arg)...)
    }

    return 
}



func (c AddEntryCommand) Do(data string) (out string, err error) {
    var sd StateData
	if data != "" {
		if err := json.Unmarshal([]byte(data), &sd); err != nil {
			dlog.Printf("Invalid data")
		}
	} else {
        dlog.Printf("data should not be empty")
    }

    entity := domain.Create(sd.Args.Description, sd.Args.Project, sd.Args.Tag, sd.Args.TimeEstimation)

    if err = c.Repo.Insert(&entity); err != nil {
        return
    }

    return
}

// descrption
func (c AddEntryCommand) generateDescriptionItems(sd StateData, enteredDescription string) (items []alfred.Item) {
    args := sd.Args
    args.Description = enteredDescription
    item := alfred.Item{
        Title: fmt.Sprintf("New description: %s", enteredDescription),
        Subtitle: c.subtitle(sd),
        Autocomplete: sd.Args.Tag,
        Arg: &alfred.ItemArg{
            Keyword: command.AddEntryKeyword,
            Mode: alfred.ModeTell,
            Data: alfred.Stringify(StateData{Current: next(sd.Current), Args: args}),
        },
    }
    items = append(items, item)
    return
}

// project
func (c AddEntryCommand) generateProjectItems(sd StateData, enteredArg string, projects []toggl.Project) (items []alfred.Item) {
    args := sd.Args
    for _, project := range projects {
        if enteredArg != "" {
            if !strings.Contains(project.Name, enteredArg) {
                continue
            }
        }
        args.Project = project.ID
        item := alfred.Item{
            Title: fmt.Sprintf("Project: %s", project.Name),
            Subtitle: c.subtitle(sd),
            Autocomplete: fmt.Sprintf("Project: %s", project.Name),
            Arg: &alfred.ItemArg{
                Keyword: command.AddEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(StateData{
                    Current: next(sd.Current),
                    Args: args,
                }),
            },
        }
        items = append(items, item)
    }
    return
}

// tag
func (c AddEntryCommand) generateTagItems(sd StateData, enteredArg string, tags []toggl.Tag) (items []alfred.Item) {
    args := sd.Args

    if enteredArg == "" {
        noTagItem := alfred.Item{
            Title: "No tag",
            Arg: &alfred.ItemArg{
                Keyword: command.AddEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(StateData{
                    Current: next(sd.Current),
                    Args: args,
                }),
            },
        }
        items = append(items, noTagItem)
    }

    for _, tag := range tags {
        if enteredArg != "" {
            if !strings.Contains(tag.Name, enteredArg) {
                continue
            }
        }
        args.Tag = tag.Name
        item := alfred.Item{
            Title: fmt.Sprintf("Tag: %s", tag.Name),
            Subtitle: c.subtitle(sd),
            Autocomplete: fmt.Sprintf("Tag: %s", tag.Name),
            Arg: &alfred.ItemArg{
                Keyword: command.AddEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(StateData{
                    Current: next(sd.Current),
                    Args: args,
                }),
            },
        }
        items = append(items, item)
    }
    return
}

// time estimation
func (c AddEntryCommand) generateTimeEstimationItems(sd StateData, enteredEstimationStr string) (items []alfred.Item) {
    args := sd.Args
    var estimationTime int
    var err error
    estimationTime, err = strconv.Atoi(enteredEstimationStr)
    if err != nil {
        estimationTime = 30     // TODO error handling
        dlog.Printf("Integer must be entered")
    }

    args.TimeEstimation = estimationTime
    item := alfred.Item{
        Title: fmt.Sprintf("Time estimation [min]: %d", estimationTime),
        Subtitle: c.subtitle(sd),
        Arg: &alfred.ItemArg{
            Keyword: command.AddEntryKeyword,
            Mode: alfred.ModeDo,
            Data: alfred.Stringify(StateData{
                Current: next(sd.Current),
                Args: args,
            }),
        },
    }
    items = append(items, item)
    return
}

func next(c state) state {
    switch c {
        case DescriptionEdit:
            return TimeEstimationEdit
        case ProjectEdit:
            return TagEdit
        case TagEdit:
            return DescriptionEdit
        case TimeEstimationEdit:
            return EndEdit
    }
    return EndEdit
}

func (c AddEntryCommand) subtitle(sd StateData) string {
    args := sd.Args
    projects, _ := c.Repo.GetProjects()
    projectName := "-"
    for _, p := range projects {
        if p.ID == args.Project {
            projectName = p.Name
        }
    }

    return fmt.Sprintf("Project: %s, Tag: %s, Desc: %s, Estimation: %d", projectName, args.Tag, args.Description, args.TimeEstimation)
}
