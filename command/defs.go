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

const iconDir = "./icons/"
const (
    OffIcon = iconDir + "power_off.png" 
    OnIcon = iconDir + "power_on.png"
    WarningIcon = iconDir + "warning.png"
    LateIcon = iconDir + "late.png"
    LateCheckedIcon = iconDir + "late_checked.png"
    BackIcon = iconDir + "back.png"
)
