package add

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"toggl_time_entry_manipulator/config"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"

	"github.com/jason0x43/go-alfred"
	"github.com/jason0x43/go-toggl"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.add]", log.LstdFlags)

type StateData struct {
    Current  state
    Entity domain.TimeEntryEntity
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

type IAddEntryCommand interface {
    alfred.Filter
    alfred.Action
}

type addEntryCommand struct {
    Repo repository.ICachedRepository
    Config *config.WorkflowConfig
}

func NewAddEntryCommand(repo repository.ICachedRepository, config *config.WorkflowConfig) (IAddEntryCommand) {
    return &addEntryCommand{Repo: repo, Config: config}
}

func (c addEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.AddEntryKeyword,
        Description: "add entry: project -> tag -> description -> estimation",
        IsEnabled: true,
    }
}

func (c addEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    // load from alfred variable
    var sd StateData

	if data != "" {
		if err := json.Unmarshal([]byte(data), &sd); err != nil {
			dlog.Printf("Invalid data")
		}
	} else {
        sd = StateData{
            Current: initialState,
            Entity: domain.TimeEntryEntity{},
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



func (c addEntryCommand) Do(data string) (out string, err error) {
    var sd StateData
	if data != "" {
		if err := json.Unmarshal([]byte(data), &sd); err != nil {
			dlog.Printf("Invalid data")
		}
	} else {
        dlog.Printf("data should not be empty")
    }

    entity := sd.Entity

    if err = c.Repo.Insert(&entity); err != nil {
        return
    }

    out = fmt.Sprintf("Time entry [%s] has started", entity.Entry.Description)

    return
}

// descrption
func (c addEntryCommand) generateDescriptionItems(sd StateData, enteredDescription string) (items []alfred.Item) {
    entity := sd.Entity
    entity.Entry.Description = enteredDescription
    tag := ""
    subtitle := ""
    if len(entity.Entry.Tags) > 0 {
        tag = entity.Entry.Tags[0]
        subtitle = fmt.Sprintf("autocomplete: %s", tag)
    }
    nextState, mode := c.next(sd.Current)
    item := alfred.Item{
        Title: fmt.Sprintf("New description: %s", enteredDescription),
        Subtitle: subtitle,
        Autocomplete: tag,
        Arg: &alfred.ItemArg{
            Keyword: command.AddEntryKeyword,
            Mode: mode,
            Data: alfred.Stringify(StateData{Current: nextState, Entity: entity}),
        },
    }
    items = append(items, item)
    if hasPrevState(sd.Current) {
        items = append(items, c.generateBackItem(sd))
    }
    return
}

// project
func (c addEntryCommand) generateProjectItems(sd StateData, enteredArg string, projects []toggl.Project) (items []alfred.Item) {
    entity := sd.Entity
    nextState, mode := c.next(sd.Current)
    items = command.GenerateItemsForProject(
        projects,
        enteredArg,
        entity,
        *c.Config,
        func(e domain.TimeEntryEntity) (alfred.ItemArg) {
            return alfred.ItemArg{
                Keyword: command.AddEntryKeyword,
                Mode: mode,
                Data: alfred.Stringify(StateData{
                    Current: nextState,
                    Entity: e,
                })}
        },
    )
    if hasPrevState(sd.Current) {
        items = append(items, c.generateBackItem(sd))
    }
    return
}

// tag
func (c addEntryCommand) generateTagItems(sd StateData, enteredArg string, tags []toggl.Tag) (items []alfred.Item) {
    entity := sd.Entity

    nextState, mode := c.next(sd.Current)
    items = command.GenerateItemsForTag(
        tags,
        enteredArg,
        entity,
        *c.Config,
        func(e domain.TimeEntryEntity) (alfred.ItemArg) {
            return alfred.ItemArg{
                Keyword: command.AddEntryKeyword,
                Mode: mode,
                Data: alfred.Stringify(StateData{
                    Current: nextState,
                    Entity: e,
                })}})
    if hasPrevState(sd.Current) {
        items = append(items, c.generateBackItem(sd))
    }
    return
}

// time estimation
func (c addEntryCommand) generateTimeEstimationItems(sd StateData, enteredEstimationStr string) (items []alfred.Item) {
    entity := sd.Entity
    var estimationTime int
    var err error
    estimationTime, err = strconv.Atoi(enteredEstimationStr)
    if err != nil {
        estimationTime = 30     // TODO error handling
        dlog.Printf("Integer must be entered")
    }

    entity.Estimation.Duration = estimationTime
    nextState, mode := c.next(sd.Current)
    item := alfred.Item{
        Title: fmt.Sprintf("Time estimation [min]: %d", estimationTime),
        Arg: &alfred.ItemArg{
            Keyword: command.AddEntryKeyword,
            Mode: mode,
            Data: alfred.Stringify(StateData{
                Current: nextState,
                Entity: entity,
            }),
        },
    }
    items = append(items, item)
    if hasPrevState(sd.Current) {
        items = append(items, c.generateBackItem(sd))
    }
    return
}

var processOrders = []state{ProjectEdit, TagEdit, DescriptionEdit, TimeEstimationEdit, EndEdit}
var processOrdersWithoutEstimation = []state{ProjectEdit, TagEdit, DescriptionEdit, EndEdit}
func (command addEntryCommand) next(c state) (state, alfred.ModeType) {
    var orders []state
    if command.Config != nil && command.Config.RecordEstimate {
        orders = processOrders
    } else {
        orders = processOrdersWithoutEstimation
    }
    n := len(orders)

    next_i := n - 1
    for i, s := range orders[:n-1] {
        if (s == c) {
            next_i = i + 1
            break
        }
    }
    var mode alfred.ModeType
    nextState := orders[next_i]
    if nextState == EndEdit {
        mode = alfred.ModeDo
    } else {
        mode = alfred.ModeTell
    }
    return nextState, mode
}

func (command addEntryCommand) prev(c state) (state, alfred.ModeType) {
    var orders []state
    if command.Config.RecordEstimate {
        orders = processOrders
    } else {
        orders = processOrdersWithoutEstimation
    }
    prev_i := 0
    for i, s := range orders[1:] {
        if (s == c) {
            prev_i = i
            break
        }
    }
    return orders[prev_i], alfred.ModeTell
}

func getPrevEntity(entity domain.TimeEntryEntity, prevState state) (prevEntity domain.TimeEntryEntity) {
    prevEntity = entity.Copy()
    switch prevState {
    case ProjectEdit:
        prevEntity.Entry.Pid = 0
    case TagEdit:
        prevEntity.Entry.Tags = []string{}
    case DescriptionEdit:
        prevEntity.Entry.Description = ""
    case TimeEstimationEdit:
        prevEntity.Estimation.Duration = 0
    }
    return prevEntity
}

func hasPrevState(c state) bool {
    return processOrders[0] != c
}

func (c addEntryCommand)generateBackItem(stateData StateData) (alfred.Item) {
    prevState, _ := c.prev(stateData.Current)
    return command.GenerateBackItem(command.AddEntryKeyword, alfred.Stringify(StateData{
        Current: prevState,
        Entity: getPrevEntity(stateData.Entity, prevState),
    }))
}
