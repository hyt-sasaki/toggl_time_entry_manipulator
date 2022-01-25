package modify

import (
	"encoding/json"
    "fmt"
	"log"
	"os"
    "strconv"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/command"
	"github.com/jason0x43/go-alfred"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.get]", log.LstdFlags)

type ModifyEntryCommand struct {
    Repo repository.ICachedRepository
}

func (c ModifyEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.ModifyEntryKeyword,
        Description: "modify entry",
        IsEnabled: true,
    }
}

func (c ModifyEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    var modifyData command.ModifyData

    err = json.Unmarshal([]byte(data), &modifyData)
    if err != nil {
        dlog.Printf("Invalid data")
        return
    }
    target := modifyData.Target

    id := modifyData.Ref.ID
    entity, err := c.Repo.FindOneById(id)
    if err != nil {
        dlog.Printf("Not found: id = %d", id)
        return
    }

    switch target {
        case command.ModifyDescription:
            entity.Entry.Description = arg
            items = append(items, alfred.Item{
                Title: fmt.Sprintf("Description: %s", arg),
                Subtitle: "Enter new description",
                Arg: &alfred.ItemArg{
                    Keyword: command.ModifyEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: alfred.Stringify(entity),
                },
            })
        case command.ModifyDuration:
            estimatedDuration, err := strconv.Atoi(arg)
            if err != nil {
                estimatedDuration = entity.Estimation.Duration
                dlog.Printf("Integer must be entered")
            }
            entity.Estimation.Duration = estimatedDuration
            items = append(items, alfred.Item{
                Title: fmt.Sprintf("Duration: %d", estimatedDuration),
                Subtitle: "Enter estimated duration",
                Arg: &alfred.ItemArg{
                    Keyword: command.ModifyEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: alfred.Stringify(entity),
                },
            })
        case command.ModifyMemo:
            entity.Estimation.Memo = arg
            items = append(items, alfred.Item{
                Title: fmt.Sprintf("Memo: %s", arg),
                Subtitle: "Enter memo",
                Arg: &alfred.ItemArg{
                    Keyword: command.ModifyEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: alfred.Stringify(entity),
                },
            })
    }

    return
}

func (c ModifyEntryCommand) Do(data string) (out string, err error) {
    var entity domain.TimeEntryEntity

    err = json.Unmarshal([]byte(data), &entity)
    if err != nil {
        dlog.Printf("Invalid data")
        return
    }

    err = c.Repo.Update(&entity)
    if err != nil {
        dlog.Printf("Failed to update entity")
        return
    }

    out = "Time entry has been updated successfully"
    return
}
