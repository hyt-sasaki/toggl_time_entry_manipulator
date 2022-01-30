package modify

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/domain"
	"toggl_time_entry_manipulator/repository"

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
            var itemArg *alfred.ItemArg = nil
            icon := "power_off.png"
            if arg != "" {
                entity.Entry.Description = arg
                itemArg = &alfred.ItemArg{
                    Keyword: command.ModifyEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: alfred.Stringify(entity)}
                icon = "power_on.png"
            }
            items = append(items, alfred.Item{
                Title: fmt.Sprintf("Description: %s", arg),
                Subtitle: "Enter new description",
                Autocomplete: entity.Entry.Description,
                Icon: icon,
                Arg: itemArg})
            items = append(items, generateBackItem(modifyData))

        case command.ModifyDuration:
            estimatedDuration, err := strconv.Atoi(arg)
            icon := "power_off.png"
            var title string
            var itemArg *alfred.ItemArg
            if err != nil {
                title = "Duration: -"
                itemArg = nil
            } else {
                entity.Estimation.Duration = estimatedDuration
                itemArg = &alfred.ItemArg{
                    Keyword: command.ModifyEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: alfred.Stringify(entity),
                }
                icon = "power_on.png"
                title = fmt.Sprintf("Duration: %d", estimatedDuration)
            }
            items = append(items, alfred.Item{
                Title: title,
                Subtitle: "Enter estimated duration",
                Autocomplete: fmt.Sprintf("%d", estimatedDuration),
                Icon: icon,
                Arg: itemArg,
            })
            items = append(items, generateBackItem(modifyData))

        case command.ModifyProject:
            projects, _ := c.Repo.GetProjects()     // TODO error handling
            items = command.GenerateItemsForProject(
                projects, arg, entity,
                func (e domain.TimeEntryEntity) (alfred.ItemArg) {
                    return alfred.ItemArg{
                        Keyword: command.ModifyEntryKeyword,
                        Mode: alfred.ModeDo,
                        Data: alfred.Stringify(e)}})
            items = append(items, generateBackItem(modifyData))

        case command.ModifyTag:
            tags, _ := c.Repo.GetTags()     // TODO error handling
            items = command.GenerateItemsForTag(
                tags, arg, entity,
                func(e domain.TimeEntryEntity) (alfred.ItemArg) {
                     return alfred.ItemArg{
                         Keyword: command.ModifyEntryKeyword,
                         Mode: alfred.ModeDo,
                         Data: alfred.Stringify(e)}})
            items = append(items, generateBackItem(modifyData))

        case command.ModifyStart:
            start, err := convertToTime(arg, entity.Entry.Start)
            autocomplete := c.calcLatestStop(entity)
            var itemArg *alfred.ItemArg
            var title string
            icon := "power_off.png"
            beforeUpdated := *entity.Entry.Start
            if err == nil {
                entity.Entry.Start = &start
                itemArg = &alfred.ItemArg{
                    Keyword: command.ModifyEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: alfred.Stringify(entity),
                }
                title = fmt.Sprintf("Start: %s", start.Format("06/01/02 15:04"))
                icon = "power_on.png"
            } else {
                itemArg = nil
                title = "Start: -"
            }

            items = append(items, alfred.Item{
                Title: title,
                Autocomplete: autocomplete,
                Subtitle: fmt.Sprintf("Modify start time (%s)", beforeUpdated.In(time.Local).Format("06/01/02 15:04")),
                Icon: icon,
                Arg: itemArg,
            })
            items = append(items, generateBackItem(modifyData))

        case command.ModifyStop:
            stop, err := convertToTime(arg, entity.Entry.Stop)
            var itemArg *alfred.ItemArg
            var title string
            icon := "poewr_off.png"
            beforeUpdated := *entity.Entry.Stop
            if err == nil {
                entity.Entry.Stop = &stop
                itemArg = &alfred.ItemArg{
                    Keyword: command.ModifyEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: alfred.Stringify(entity),
                }
                title = fmt.Sprintf("Stop: %s", stop.Format("06/01/02 15:04"))
                icon = "power_on.png"
            } else {
                itemArg = nil
                title = "Stop: -"
            }

            layout := "06/01/02 15:04"
            items = append(items, alfred.Item{
                Title: title,
                Subtitle: fmt.Sprintf("Modify stop time (%s)", beforeUpdated.In(time.Local).Format(layout)),
                Autocomplete: entity.Entry.Stop.In(time.Local).Format(layout),
                Icon: icon,
                Arg: itemArg,
            })
            items = append(items, generateBackItem(modifyData))

        case command.ModifyMemo:
            var itemArg *alfred.ItemArg = nil
            icon := "poewr_off.png"
            if arg != "" {
                entity.Estimation.Memo = arg
                itemArg = &alfred.ItemArg{
                    Keyword: command.ModifyEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: alfred.Stringify(entity) }
                icon = "power_on.png"
            }
            items = append(items, alfred.Item{
                Title: fmt.Sprintf("Memo: %s", arg),
                Subtitle: "Enter memo",
                Autocomplete: entity.Estimation.Memo,
                Icon: icon,
                Arg: itemArg,
            })
            items = append(items, generateBackItem(modifyData))
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

func convertToTime(dateStr string, base *time.Time) (result time.Time, err error) {
    layout := "06/01/02 15:04"
    date := base.In(time.Local).Format("06/01/02")
    fullDateStr := fmt.Sprintf("%s %s", date, dateStr)
    result, fail := time.ParseInLocation(layout, fullDateStr, time.Local)

    if fail == nil {
        return
    }

    result, err = time.ParseInLocation(layout, dateStr, time.Local)

    return
}

func (c ModifyEntryCommand) calcLatestStop(entity domain.TimeEntryEntity) (out string) {
    layout := "06/01/02 15:04"
    // 修正対象のentityがstopしている場合はentityの値をそのまま使う
    if !entity.IsRunning() {
        out = entity.Entry.Start.In(time.Local).Format(layout)
        return
    }

    entities, _ := c.Repo.Fetch()   // sort済み
    // entityが1個しかない場合は何もしない
    if len(entities) < 2 {
        return
    }
    // entityが最新のものではない場合は何もしない
    if (entities[0].Entry.ID != entity.Entry.ID) {
        return
    }
    latestStop := entities[1].Entry.Stop
    out = latestStop.In(time.Local).Format(layout)
    return
}

func generateBackItem(modifyData command.ModifyData) (alfred.Item) {
    return command.GenerateBackItem(command.GetEntryKeyword, alfred.Stringify(command.DetailRefData{
        ID: modifyData.Ref.ID,
    }))
}
