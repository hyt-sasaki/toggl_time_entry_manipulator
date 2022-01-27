package add

import (
	"encoding/json"
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

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.add]", log.LstdFlags)

type AddEntryCommand struct {
    Repo repository.ICachedRepository
}

type stateData struct {
    Current  state
    Args entryArgs
}
type state int
const (
    DescriptionEdit state = iota
    ProjectEdit
    TagEdit 
    TimeEstimationEdit
    EndEdit
)
const initialState = DescriptionEdit
type entryArgs struct {
    Description string
    Project int
    Tag string
    TimeEstimation int  // minutes
}

func (c AddEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.AddEntryKeyword,
        Description: "add entry: description -> project -> tag",
        IsEnabled: true,
    }
}

func (c AddEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    // load from alfred variable
    var sd stateData

	if data != "" {
		if err := json.Unmarshal([]byte(data), &sd); err != nil {
			dlog.Printf("Invalid data")
		}
	} else {
        sd = stateData{
            Current: initialState,
        }
    }


    // generate items
    switch sd.Current {
        case DescriptionEdit:
            items = append(items, generateDescriptionItems(sd, arg)...)
        case ProjectEdit:
            var projects []toggl.Project
            projects, err = c.Repo.GetProjects()
            if err != nil {
                return
            }
            items = append(items, generateProjectItems(sd, arg, projects)...)
        case TagEdit:
            var tags []toggl.Tag
            tags, err = c.Repo.GetTags()
            if err != nil {
                return
            }
            items = append(items, generateTagItems(sd, arg, tags)...)
        case TimeEstimationEdit:
            items = append(items, generateTimeEstimationItems(sd, arg)...)
    }

    return 
}



func (c AddEntryCommand) Do(data string) (out string, err error) {
    var sd stateData
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
func generateDescriptionItems(sd stateData, enteredDescription string) (items []alfred.Item) {
    args := sd.Args
    args.Description = enteredDescription
    item := alfred.Item{
        Title: fmt.Sprintf("New description: %s", enteredDescription),
        Subtitle: "Create new description for your time entry",
        Arg: &alfred.ItemArg{
            Keyword: command.AddEntryKeyword,
            Mode: alfred.ModeTell,
            Data: alfred.Stringify(stateData{Current: next(sd.Current), Args: args}),
        },
    }
    items = append(items, item)
    return
}

// project
func generateProjectItems(sd stateData, enteredArg string, projects []toggl.Project) (items []alfred.Item) {
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
            Subtitle: "Select the project for your time entry",
            Autocomplete: fmt.Sprintf("Project: %s", project.Name),
            Arg: &alfred.ItemArg{
                Keyword: command.AddEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(stateData{
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
func generateTagItems(sd stateData, enteredArg string, tags []toggl.Tag) (items []alfred.Item) {
    args := sd.Args

    if enteredArg == "" {
        noTagItem := alfred.Item{
            Title: "No tag",
            Arg: &alfred.ItemArg{
                Keyword: command.AddEntryKeyword,
                Mode: alfred.ModeDo,
                Data: alfred.Stringify(stateData{
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
            Subtitle: "Select the tag for your time entry",
            Autocomplete: fmt.Sprintf("Tag: %s", tag.Name),
            Arg: &alfred.ItemArg{
                Keyword: command.AddEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(stateData{
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
func generateTimeEstimationItems(sd stateData, enteredEstimationStr string) (items []alfred.Item) {
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
        Subtitle: "Enter time estimation for your time entry (default: 30 min)",
        Arg: &alfred.ItemArg{
            Keyword: command.AddEntryKeyword,
            Mode: alfred.ModeDo,
            Data: alfred.Stringify(stateData{
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
            return ProjectEdit
        case ProjectEdit:
            return TagEdit
        case TagEdit:
            return TimeEstimationEdit
        case TimeEstimationEdit:
            return EndEdit
    }
    return EndEdit
}
