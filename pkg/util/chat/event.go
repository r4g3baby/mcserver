package chat

type (
	ClickAction string
	HoverAction string

	ClickEvent struct {
		Action ClickAction
		Value  string
	}

	HoverEvent struct {
		Action   HoverAction
		Contents interface{}
	}
)

const (
	OpenUrlClickAction        ClickAction = "open_url"
	OpenFileClickAction       ClickAction = "open_file"
	RunCommandClickAction     ClickAction = "run_command"
	SuggestCommandClickAction ClickAction = "suggest_command"
	ChangePageClickAction     ClickAction = "change_page"
	CopyToClipboardAction     ClickAction = "copy_to_clipboard"

	ShowTextHoverAction   HoverAction = "show_text"
	ShowItemHoverAction   HoverAction = "show_item"
	ShowEntityHoverAction HoverAction = "show_entity"
)
