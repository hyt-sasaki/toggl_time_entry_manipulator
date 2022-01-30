package tests

import (
    "encoding/json"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/domain"
)

func StringifyDetailRefData(data command.DetailRefData) (string) {
    dataBytes, _ := json.Marshal(data)
    return string(dataBytes)
}

func StringifyEntity(entity domain.TimeEntryEntity) (string) {
    dataBytes, _ := json.Marshal(entity)
    return string(dataBytes)
}
