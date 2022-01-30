// +build wireinject

package main

import (
	"github.com/jason0x43/go-alfred"
    "github.com/google/wire"
    "toggl_time_entry_manipulator/client"
    "toggl_time_entry_manipulator/config"
    "toggl_time_entry_manipulator/repository"
    "toggl_time_entry_manipulator/repository/myCache"
    "toggl_time_entry_manipulator/command/add"
    "toggl_time_entry_manipulator/command/list"
    "toggl_time_entry_manipulator/command/get"
    "toggl_time_entry_manipulator/command/modify"
    "toggl_time_entry_manipulator/command/stop"
    "toggl_time_entry_manipulator/command/delete"
    "toggl_time_entry_manipulator/command/continue_entry"
)

func initializeCommands(workflow alfred.Workflow, firstCall bool) (commands []alfred.Command, err error) {
    wire.Build(
        NewServiceAccount,
        NewCacheFile,
        NewCache,
        NewConfigFile,
        NewConfig,
        NewCommands,
        client.NewEstimationClient,
        client.NewTogglClient,
        repository.NewTimeEntryRepository,
        repository.NewCachedRepository,
        add.NewAddEntryCommand,
        list.NewListEntryCommand,
        get.NewGetEntryCommand,
        modify.NewModifyEntryCommand,
        stop.NewStopEntryCommand,
        delete.NewDeleteEntryCommand,
        continue_entry.NewContinueEntryCommand,
        wire.FieldsOf(new(*config.Config), "TogglConfig"),
        wire.FieldsOf(new(*config.Config), "FirestoreConfig"),
        wire.FieldsOf(new(*config.Config), "WorkflowConfig"),
        wire.Bind(new(client.ITogglClient), new(*client.TogglClient)),
        wire.Bind(new(client.IEstimationClient), new(*client.EstimationClient)),
        wire.Bind(new(repository.ITimeEntryRepository), new(*repository.TimeEntryRepository)),
        wire.Bind(new(myCache.ICache), new(*myCache.Cache)),
        wire.Bind(new(repository.ICachedRepository), new(*repository.CachedRepository)),
    )
    return []alfred.Command{}, nil
}
