package main

import (
    "fmt"
    "strings"
    "encoding/json"
    "strconv"
    "time"

    "github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
    "toggl_time_entry_manipulator/estimation_client"
)

type AddEntryCommand struct {
    firestoreClient *estimation_client.EstimationClient
}

const AddEntryKeyword = "add_entry"
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
type entryArgs struct {
    Description string
    Project int
    Tag string
    TimeEstimation int  // minutes
}

func (c AddEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: AddEntryKeyword,
        Description: "add entry: description -> project -> tag",
        IsEnabled: true,
    }
}

func (c AddEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    // fetch toggl info
	if err = checkRefresh(); err != nil {
		return
	}

    // load from alfred variable
    var sd stateData

	if data != "" {
		if err := json.Unmarshal([]byte(data), &sd); err != nil {
			dlog.Printf("Invalid data")
		}
	} else {
        sd = stateData{
            Current: DescriptionEdit,
        }
    }
    dlog.Printf("sd data is %s", data)


    // generate items
    switch sd.Current {
        case DescriptionEdit:
            items = append(items, generateDescriptionItems(sd.Args, arg)...)
        case ProjectEdit:
            items = append(items, generateProjectItems(sd.Args, arg)...)
        case TagEdit:
            items = append(items, generateTagItems(sd.Args, arg)...)
        case TimeEstimationEdit:
            items = append(items, generateTimeEstimationItems(sd.Args, arg)...)
    }

    return 
}



func (c AddEntryCommand) Do(data string) (out string, err error) {
    dlog.Printf("data is %s", data)
	session := toggl.OpenSession(config.APIKey)

    var sd stateData
	if data != "" {
		if err := json.Unmarshal([]byte(data), &sd); err != nil {
			dlog.Printf("Invalid data")
		}
	} else {
        dlog.Printf("data should not be empty")
    }

    var time_entry toggl.TimeEntry
    time_entry, err = session.StartTimeEntryForProject(sd.Args.Description, sd.Args.Project, false)
    if err != nil {
        dlog.Printf("Something wrong. %s", err)
    }
    if sd.Args.Tag != "" {
        time_entry.Tags = []string{sd.Args.Tag}
        _, err = session.UpdateTimeEntry(time_entry)
        if err != nil {
            dlog.Printf("Something wrong. %s", err)
        }
    }

    timeEstimation := sd.Args.TimeEstimation
    if timeEstimation != 0 {
        dlog.Printf("TimeEstimation is %d", timeEstimation)

        if err = c.firestoreClient.Insert(strconv.Itoa(time_entry.ID), estimation_client.Estimation{
            Duration: timeEstimation,
            Memo: "",
            CreatedTm: time.Now(),
            UpdatedTm: time.Now(),
        }); err != nil {
            dlog.Printf("Failed to add time entry estimation: %v", err)
        }

    } else {
        dlog.Printf("TimeEstimation is not entered")
    }

    return
}

// descrption
func generateDescriptionItems(args entryArgs, enteredDescription string) (items []alfred.Item) {
    item := alfred.Item{
        Title: fmt.Sprintf("New description: %s", enteredDescription),
        Subtitle: "Create new description for your time entry",
        Arg: &alfred.ItemArg{
            Keyword: AddEntryKeyword,
            Mode: alfred.ModeTell,
            Data: alfred.Stringify(stateData{Current: ProjectEdit, Args: entryArgs{Description: enteredDescription}}),
        },
    }
    items = append(items, item)
    return
}

// project
func generateProjectItems(args entryArgs, enteredArg string) (items []alfred.Item) {
    projects := cache.Account.Data.Projects
    for _, project := range projects {
        if enteredArg != "" {
            if !strings.Contains(project.Name, enteredArg) {
                continue
            }
        }
        item := alfred.Item{
            Title: fmt.Sprintf("Project: %s", project.Name),
            Subtitle: "Select the project for your time entry",
            Autocomplete: fmt.Sprintf("Project: %s", project.Name),
            Arg: &alfred.ItemArg{
                Keyword: AddEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(stateData{
                    Current: TagEdit,
                    Args: generateArgsOfProjectItem(args, project),
                }),
            },
        }
        items = append(items, item)
    }
    return
}

func generateArgsOfProjectItem(args entryArgs, project toggl.Project) (out entryArgs) {
    out = entryArgs{
        Description: args.Description,
        Project: project.ID,
        Tag: args.Tag,
    }
    return
}

// tag
func generateTagItems(args entryArgs, enteredArg string) (items []alfred.Item) {
    tags := cache.Account.Data.Tags

    if enteredArg == "" {
        noTagItem := alfred.Item{
            Title: "No tag",
            Arg: &alfred.ItemArg{
                Keyword: AddEntryKeyword,
                Mode: alfred.ModeDo,
                Data: alfred.Stringify(stateData{
                    Current: TimeEstimationEdit,
                    Args: generateArgsOfTagItem(args, toggl.Tag{}),
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
        item := alfred.Item{
            Title: fmt.Sprintf("Tag: %s", tag.Name),
            Subtitle: "Select the tag for your time entry",
            Autocomplete: fmt.Sprintf("Tag: %s", tag.Name),
            Arg: &alfred.ItemArg{
                Keyword: AddEntryKeyword,
                Mode: alfred.ModeTell,
                Data: alfred.Stringify(stateData{
                    Current: TimeEstimationEdit,
                    Args: generateArgsOfTagItem(args, tag),
                }),
            },
        }
        items = append(items, item)
    }
    return
}

func generateArgsOfTagItem(args entryArgs, tag toggl.Tag) (out entryArgs) {
    out = entryArgs{
        Description: args.Description,
        Project: args.Project,
        Tag: tag.Name,
    }
    return
}

// time estimation
func generateTimeEstimationItems(args entryArgs, enteredEstimationStr string) (items []alfred.Item) {
    var estimationTime int
    var err error
    estimationTime, err = strconv.Atoi(enteredEstimationStr)
    if err != nil {
        estimationTime = 30     // TODO error handling
        dlog.Printf("Integer must be entered")
    }

    item := alfred.Item{
        Title: fmt.Sprintf("Time estimatune [min]: %d", estimationTime),
        Subtitle: "Enter time estimation for your time entry (default: 30 min)",
        Arg: &alfred.ItemArg{
            Keyword: AddEntryKeyword,
            Mode: alfred.ModeDo,
            Data: alfred.Stringify(stateData{
                Current: EndEdit,
                Args: generateArgsOfTimeEstimationItem(args, estimationTime),
            }),
        },
    }
    items = append(items, item)
    return
}

func generateArgsOfTimeEstimationItem(args entryArgs, estimationTime int) (out entryArgs) {
    out = entryArgs{
        Description: args.Description,
        Project: args.Project,
        Tag: args.Tag,
        TimeEstimation: estimationTime,
    }
    return
}
