package tests

import (
    "encoding/json"
	"toggl_time_entry_manipulator/command"
)

func StringifyDetailRefData(data command.DetailRefData) (string) {
    dataBytes, _ := json.Marshal(data)
    return string(dataBytes)
}
