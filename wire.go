// +build wireinject

package main

import (
	"github.com/jason0x43/go-alfred"
    "github.com/google/wire"
    "toggl_time_entry_manipulator/estimation_client"
    "toggl_time_entry_manipulator/repository"
    "toggl_time_entry_manipulator/repository/myCache"
	"google.golang.org/api/option"
)

func initializeRepository(workflow alfred.Workflow, serviceAccount option.ClientOption) (repo *repository.CachedRepository, err error) {
    wire.Build(
        NewCacheFile,
        NewCache,
        NewConfigFile,
        NewConfig,
        estimation_client.NewEstimationClient,
        estimation_client.NewTogglClient,
        repository.NewTimeEntryRepository,
        repository.NewCachedRepository,
        wire.FieldsOf(new(*repository.Config), "TogglAPIKey"),
        wire.Bind(new(estimation_client.ITogglClient), new(*estimation_client.TogglClient)),
        wire.Bind(new(estimation_client.IEstimationClient), new(*estimation_client.EstimationClient)),
        wire.Bind(new(myCache.ICache), new(*myCache.Cache)),
    )
    return &repository.CachedRepository{}, nil
}
