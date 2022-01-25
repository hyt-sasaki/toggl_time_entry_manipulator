package command

const AddEntryKeyword = "add_entry"
const StopEntryKeyword = "stop_entry"
const GetEntryKeyword = "get_entry"
const ModifyEntryKeyword = "modify_entry"
const ListEntryKeyword = "list_entries"

type DetailRefData struct {
    ID int
}

type ModifyData struct {
    Ref DetailRefData
    Target modifyTarget
}

type modifyTarget int
const (
    ModifyDescription modifyTarget = iota
    ModifyDuration
    ModifyStart 
    ModifyStop
    ModifyMemo
)
