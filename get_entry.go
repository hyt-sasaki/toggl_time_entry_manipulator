package main

import (
    "fmt"
    "strconv"
    "github.com/jason0x43/go-alfred"
    "github.com/jason0x43/go-toggl"
    "toggl_time_entry_manipulator/estimation_client"
)

type GetEntryCommand struct {
    firestoreClient estimation_client.IEstimationClient
}

const GetEntryKeyword = "get_entry"
type EntryWithEstimation struct {
    Entry toggl.TimeEntry
    Estimation estimation_client.Estimation
}

func (c GetEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: GetEntryKeyword,
        Description: "get entries",
        IsEnabled: true,
    }
}

func (c GetEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    dlog.Printf("Items")
    entryWithEstimations, err := c.getEntries()
    for _, entryWithEstimation := range entryWithEstimations {
        item := alfred.Item{
            Title: fmt.Sprintf("Description: %s", entryWithEstimation.Entry.Description),
            Subtitle: fmt.Sprintf("actual duration: %s, estimation: %d", convertDuration(entryWithEstimation.Entry.Duration), entryWithEstimation.Estimation.Duration),
            Arg: &alfred.ItemArg{
                Keyword: GetEntryKeyword,
                Mode: alfred.ModeTell,
            },
        }
        items = append(items, item)
    }
    return
}

func (c GetEntryCommand) getEntries() (entryWithEstimations []EntryWithEstimation, err error){
    // fetch toggl info
	if err = checkRefresh(); err != nil {
		return
	}
    dlog.Printf("getEntries")
    entries := cache.Account.Data.TimeEntries
    dlog.Println(entries)
    var entryIds []int64
    for _, entry := range entries {
        entryIds = append(entryIds, int64(entry.ID))
    }
    estimations, err := c.firestoreClient.Fetch(entryIds)
    dlog.Printf("estimations")
    for idx, estimation := range estimations {
        entryWithEstimations = append(entryWithEstimations, EntryWithEstimation{
            Entry: entries[idx],
            Estimation: estimation,
        })
    }
    
    return 
}

func convertDuration(duration int64) string {
    if duration < 0 {
        return "[stil running...]"
    }
    min := int(duration / 60)
    return strconv.Itoa(min)
}
