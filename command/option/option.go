package option

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
    "errors"
	"toggl_time_entry_manipulator/command"
	"toggl_time_entry_manipulator/config"
	"toggl_time_entry_manipulator/repository"

	"github.com/jason0x43/go-alfred"
)

var dlog = log.New(os.Stderr, "[toggl_time_entry_manipulator.command.option]", log.LstdFlags)


type IOptionCommand interface {
    alfred.Filter
    alfred.Action
}

type OptionCommand struct {
    repo repository.ICachedRepository
    config *config.Config
    configFile config.ConfigFile
}

func NewOptionCommand(repo repository.ICachedRepository, config *config.Config, configFile config.ConfigFile) (com IOptionCommand, err error) {
    if config == nil {
        err = errors.New("Workflow config is nil.")
        return
    }
    com = OptionCommand{repo: repo, config: config, configFile: configFile}
    return
}

func (c OptionCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: command.OptionKeyword,
        Description: "option",
        IsEnabled: true,
    }
}


func (c OptionCommand) Items(arg, data string) (items []alfred.Item, err error) {
    copied := *c.config
    copied.TogglConfig.APIKey = arg
    item := alfred.Item{
        Title: fmt.Sprintf("Toggl API key: %s", arg),
        Subtitle: fmt.Sprintf("Old API key: %s", c.config.TogglConfig.APIKey),
        Autocomplete: c.config.TogglConfig.APIKey,
        Arg: &alfred.ItemArg{
            Keyword: command.OptionKeyword,
            Mode: alfred.ModeDo,
            Data: alfred.Stringify(copied),
        },
    }
    items = append(items, item)
    return
}

func (c OptionCommand) Do(data string) (out string, err error) {
    var newConfig config.Config
	if data != "" {
		if err := json.Unmarshal([]byte(data), &newConfig); err != nil {
			dlog.Printf("Invalid data")
		}
	} else {
        dlog.Printf("data should not be empty")
    }

    c.config = &newConfig
    alfred.SaveJSON(string(c.configFile), *c.config)
    out = "Config has been saved"
    return
}
