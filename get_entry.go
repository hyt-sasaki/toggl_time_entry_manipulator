package main

import (
    "github.com/jason0x43/go-alfred"
)

type GetEntryCommand struct {}
const GetEntryKeyword = "get_entry"

func (c GetEntryCommand) About() alfred.CommandDef {
    return alfred.CommandDef{
        Keyword: GetEntryKeyword,
        Description: "get entries",
        IsEnabled: true,
    }
}

func (c GetEntryCommand) Items(arg, data string) (items []alfred.Item, err error) {
    return
}
