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
            var itemArg *alfred.ItemArg
            if err != nil {
                estimatedDuration = entity.Estimation.Duration
                dlog.Printf("Integer must be entered")
                itemArg = nil
            } else {
                itemArg = &alfred.ItemArg{
                    Keyword: command.ModifyEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: alfred.Stringify(entity),
                }
            }
            entity.Estimation.Duration = estimatedDuration
            items = append(items, alfred.Item{
                Title: fmt.Sprintf("Duration: %d", estimatedDuration),
                Subtitle: "Enter estimated duration",
                Arg: itemArg,
            })

        case command.ModifyProject:
            projects, _ := c.Repo.GetProjects()     // TODO error handling
            items = command.GenerateItemsForProject(
                projects, arg, entity,
                func (e domain.TimeEntryEntity) (alfred.ItemArg) {
                    return alfred.ItemArg{
                        Keyword: command.ModifyEntryKeyword,
                        Mode: alfred.ModeDo,
                        Data: alfred.Stringify(e)}})

        case command.ModifyTag:
            tags, _ := c.Repo.GetTags()     // TODO error handling
            items = command.GenerateItemsForTag(
                tags, arg, entity,
                func(e domain.TimeEntryEntity) (alfred.ItemArg) {
                     return alfred.ItemArg{
                         Keyword: command.ModifyEntryKeyword,
                         Mode: alfred.ModeDo,
                         Data: alfred.Stringify(e)}})

        case command.ModifyStart:
            start, err := convertToTime(arg, entity.Entry.Start)
            autocomplete := c.calcLatestStop(entity)
            var itemArg *alfred.ItemArg
            var title string
            beforeUpdated := *entity.Entry.Start
            if err == nil {
                entity.Entry.Start = &start
                itemArg = &alfred.ItemArg{
                    Keyword: command.ModifyEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: alfred.Stringify(entity),
                }
                title = fmt.Sprintf("Start: %s", start.Format("06/01/02 15:04"))
            } else {
                itemArg = nil
                title = "Start: -"
            }

            items = append(items, alfred.Item{
                Title: title,
                Autocomplete: autocomplete,
                Subtitle: fmt.Sprintf("Modify start time (%s)", beforeUpdated.In(time.Local).Format("06/01/02 15:04")),
                Arg: itemArg,
            })
        case command.ModifyStop:
            stop, err := convertToTime(arg, entity.Entry.Stop)
            var itemArg *alfred.ItemArg
            var title string
            beforeUpdated := *entity.Entry.Stop
            if err == nil {
                entity.Entry.Stop = &stop
                itemArg = &alfred.ItemArg{
                    Keyword: command.ModifyEntryKeyword,
                    Mode: alfred.ModeDo,
                    Data: alfred.Stringify(entity),
                }
                title = fmt.Sprintf("Stop: %s", stop.Format("06/01/02 15:04"))
            } else {
                itemArg = nil
                title = "Stop: -"
            }

            items = append(items, alfred.Item{
                Title: title,
                Subtitle: fmt.Sprintf("Modify stop time (%s)", beforeUpdated.In(time.Local).Format("06/01/02 15:04")),
                Arg: itemArg,
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
    // 修正対象のentityがstopしている場合は何もしない
    if !entity.IsRunning() {
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
    out = latestStop.In(time.Local).Format("06/01/02 15:04")
    return
}
