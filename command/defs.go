package command

const AddEntryKeyword = "add_entry"
const StopEntryKeyword = "stop_entry"
const DeleteEntryKeyword = "delete_entry"
const GetEntryKeyword = "get_entry"
const ModifyEntryKeyword = "modify_entry"
const ListEntryKeyword = "list_entries"
const ContinueEntryKeyword = "continue_entry"
const OptionKeyword = "option"

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
    ModifyProject
    ModifyTag
    ModifyStart 
    ModifyStop
    ModifyMemo
)
