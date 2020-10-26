package chat

type (
	Component interface {
		SetBold(bold bool)
		IsBold() bool
		SetItalic(italic bool)
		IsItalic() bool
		SetUnderlined(underlined bool)
		IsUnderlined() bool
		SetStrikethrough(strikethrough bool)
		IsStrikethrough() bool
		SetObfuscated(obfuscated bool)
		IsObfuscated() bool
		SetColor(color *Color)
		GetColor() *Color
		SetInsertion(insertion string)
		GetInsertion() string
		SetClickEvent(clickEvent *ClickEvent)
		GetClickEvent() *ClickEvent
		SetHoverEvent(hoverEvent *HoverEvent)
		GetHoverEvent() *HoverEvent
		SetExtra(extra []Component)
		GetExtra() []Component
	}

	BaseComponent struct {
		Bold          bool
		Italic        bool
		Underlined    bool
		Strikethrough bool
		Obfuscated    bool
		Color         *Color
		Insertion     string
		ClickEvent    *ClickEvent
		HoverEvent    *HoverEvent
		Extra         []Component
	}

	TextComponent struct {
		Text string

		BaseComponent
	}

	TranslatableComponent struct {
		Translate string
		With      []Component

		BaseComponent
	}

	KeybindComponent struct {
		Keybind string

		BaseComponent
	}

	ScoreComponent struct {
		Score Score

		BaseComponent
	}

	Score struct {
		Name      string
		Objective string
		Value     string
	}

	SelectorComponent struct {
		Selector string

		BaseComponent
	}
)

func (c *BaseComponent) SetBold(bold bool) {
	c.Bold = bold
}

func (c *BaseComponent) IsBold() bool {
	return c.Bold
}

func (c *BaseComponent) SetItalic(italic bool) {
	c.Italic = italic
}

func (c *BaseComponent) IsItalic() bool {
	return c.Italic
}

func (c *BaseComponent) SetUnderlined(underlined bool) {
	c.Underlined = underlined
}

func (c *BaseComponent) IsUnderlined() bool {
	return c.Underlined
}

func (c *BaseComponent) SetStrikethrough(strikethrough bool) {
	c.Strikethrough = strikethrough
}

func (c *BaseComponent) IsStrikethrough() bool {
	return c.Strikethrough
}

func (c *BaseComponent) SetObfuscated(obfuscated bool) {
	c.Obfuscated = obfuscated
}

func (c *BaseComponent) IsObfuscated() bool {
	return c.Obfuscated
}

func (c *BaseComponent) SetColor(color *Color) {
	c.Color = color
}

func (c *BaseComponent) GetColor() *Color {
	return c.Color
}

func (c *BaseComponent) SetInsertion(insertion string) {
	c.Insertion = insertion
}

func (c *BaseComponent) GetInsertion() string {
	return c.Insertion
}

func (c *BaseComponent) SetClickEvent(clickEvent *ClickEvent) {
	c.ClickEvent = clickEvent
}

func (c *BaseComponent) GetClickEvent() *ClickEvent {
	return c.ClickEvent
}

func (c *BaseComponent) SetHoverEvent(hoverEvent *HoverEvent) {
	c.HoverEvent = hoverEvent
}

func (c *BaseComponent) GetHoverEvent() *HoverEvent {
	return c.HoverEvent
}

func (c *BaseComponent) SetExtra(extra []Component) {
	c.Extra = extra
}

func (c *BaseComponent) GetExtra() []Component {
	return c.Extra
}
