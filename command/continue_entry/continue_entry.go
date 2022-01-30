package continue_entry

import (
    "fmt"
	"encoding/json"
	"log"
	"os"
    "strconv"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"
	"toggl_time_entry_manipulator/command"
	"github.com/jason0x43/go-alfred"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.continue]", log.LstdFlags)

type ContinueEntryCommand struct {
    Repo repository.ICachedRepository
}

func (c ContinueEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.ContinueEntryKeyword,
        Description: "continue entry",
        IsEnabled: true,
    }
}

func (c ContinueEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    entity, err := c.getEntityFromId(data)
    if err != nil {
        return
    }

    autocomplete := fmt.Sprintf("%d", entity.Estimation.Duration)
    estimatedDuration, parseErr := strconv.Atoi(arg)
    icon := "power_off.png"
    var title string
    var itemArg *alfred.ItemArg
    if parseErr != nil {
        title = "Duration: -"
        itemArg = nil
    } else {
        entity.Estimation.Duration = estimatedDuration
        itemArg = &alfred.ItemArg{
            Keyword: command.ContinueEntryKeyword,
            Mode: alfred.ModeDo,
            Data: alfred.Stringify(entity),
        }
        icon = "power_on.png"
        title = fmt.Sprintf("Duration: %d", estimatedDuration)
    }
    items = append(items, alfred.Item{
        Title: title,
        Subtitle: "Enter estimated duration",
        Autocomplete: autocomplete,
        Icon: icon,
        Arg: itemArg,
    })

    items = append(items, generateBackItem(entity.Entry.ID))

    return
}

func (c ContinueEntryCommand) Do(data string) (out string, err error) {
    entity, err := getEntity(data)

    newEntity, err := c.Repo.Continue(&entity)
    if err != nil {
        return
    }
    
    out = fmt.Sprintf("Entry has been copied. Description: %s", newEntity.Entry.Description)
    return
}

func (c ContinueEntryCommand) getEntityFromId(data string) (entity domain.TimeEntryEntity, err error) {
    var itemData command.DetailRefData

    err = json.Unmarshal([]byte(data), &itemData)
    if err != nil {
        dlog.Printf("Invalid data")
        return
    }

    entity, err = c.Repo.FindOneById(itemData.ID)
    if err != nil {
        return
    }
    return
}

func getEntity(data string) (entity domain.TimeEntryEntity, err error) {
    err = json.Unmarshal([]byte(data), &entity)
    if err != nil {
        dlog.Printf("Invalid data")
        return
    }
    return
}

func generateBackItem(id int) (alfred.Item) {
    return command.GenerateBackItem(command.GetEntryKeyword, alfred.Stringify(command.DetailRefData{
        ID: id,
    }))
}
