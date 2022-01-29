package addtest

import (
	_ "toggl_time_entry_manipulator/supports"
	"testing"
	"github.com/stretchr/testify/assert"

    "encoding/json"
	"github.com/jason0x43/go-alfred"

	"toggl_time_entry_manipulator/command/add"
)

func convertAddStateData(data add.StateData) (dataStr string){
    dataBytes, _ := json.Marshal(data)
    dataStr = string(dataBytes)
    return
}

func assertAddItemArg(t *testing.T, actualItem alfred.Item, expectedArg add.StateData, expectedMode alfred.ModeType) {
    actualItemArg := actualItem.Arg
    assert.Equal(t, expectedMode, actualItemArg.Mode)
    var itemData add.StateData
    json.Unmarshal([]byte(actualItemArg.Data), &itemData)
    assert.Equal(t, expectedArg, itemData)
}
