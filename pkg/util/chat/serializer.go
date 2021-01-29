package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Serializer is a serializer for Minecraft chat components.
type Serializer struct {
	// Since Minecraft 1.16+ hoverEvent "value" is deprecated in favour of "contents".
	// This setting decides whether to use the "value" key or the new "contents" key.
	//
	// It is false by default to support latest minecraft protocol.
	LegacyHoverEvent bool

	// Since Minecraft 1.16+ there can be hex colors instead of named colors.
	// This option decides whether to force the use of named legacy colors by finding
	// the nearest named color of the hex color.
	//
	// It is false by default to allow hex colors.
	ForceLegacyColors bool
}

var serializer = Serializer{}

func ToLegacyText(c []Component) string {
	return serializer.ToLegacyText(c)
}

func ToJSON(c []Component) ([]byte, error) {
	return serializer.ToJSON(c)
}

func FromJSON(data []byte) ([]Component, error) {
	return serializer.FromJSON(data)
}

func (s *Serializer) ToLegacyText(components []Component) string {
	var text strings.Builder
	for _, c := range components {
		if c.IsBold() {
			text.WriteString(ColorChar + Bold.Code)
		}
		if c.IsItalic() {
			text.WriteString(ColorChar + Italic.Code)
		}
		if c.IsObfuscated() {
			text.WriteString(ColorChar + Obfuscated.Code)
		}
		if c.IsStrikethrough() {
			text.WriteString(ColorChar + Strikethrough.Code)
		}
		if c.IsUnderlined() {
			text.WriteString(ColorChar + Underline.Code)
		}

		color := c.GetColor()
		if color != nil {
			text.WriteString(ColorChar)
			if color.Name == "" {
				text.WriteString(FindNearest(*color).Code)
			} else {
				text.WriteString(color.Code)
			}
		}

		switch t := c.(type) {
		case *TextComponent:
			text.WriteString(t.Text)
		}

		text.WriteString(s.ToLegacyText(c.GetExtra()))
	}
	return text.String()
}

func (s *Serializer) ToJSON(components []Component) ([]byte, error) {
	var array []map[string]interface{}
	for _, c := range components {
		var obj = make(map[string]interface{})
		if err := s.encode(obj, c); err != nil {
			return nil, err
		}
		array = append(array, obj)
	}
	return json.Marshal(array)
}

func (s *Serializer) FromJSON(data []byte) ([]Component, error) {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	var components []Component
	if array, ok := obj.([]interface{}); ok {
		for _, obj := range array {
			c, err := s.decode(obj)
			if err != nil {
				return nil, err
			}
			components = append(components, c)
		}
	} else {
		c, err := s.decode(obj)
		if err != nil {
			return nil, err
		}
		components = append(components, c)
	}
	return components, nil
}

func (s *Serializer) encode(obj map[string]interface{}, c Component) error {
	switch t := c.(type) {
	case *TextComponent:
		return s.encodeText(obj, t)
	case *TranslatableComponent:
		return s.encodeTranslatable(obj, t)
	case *KeybindComponent:
		return s.encodeKeybind(obj, t)
	case *ScoreComponent:
		return s.encodeScore(obj, t)
	case *SelectorComponent:
		return s.encodeSelector(obj, t)
	default:
		return fmt.Errorf("unsupported component type %T", c)
	}
}

func (s *Serializer) encodeText(obj map[string]interface{}, c *TextComponent) error {
	obj["text"] = c.Text
	return s.encodeComponent(obj, c)
}

func (s *Serializer) encodeTranslatable(obj map[string]interface{}, c *TranslatableComponent) error {
	obj["translate"] = c.Translate

	var with []map[string]interface{}
	for _, w := range c.With {
		var comp = make(map[string]interface{})
		if err := s.encode(comp, w); err != nil {
			return err
		}
		with = append(with, comp)
	}
	if len(with) > 0 {
		obj["with"] = with
	}

	return s.encodeComponent(obj, c)
}

func (s *Serializer) encodeKeybind(obj map[string]interface{}, c *KeybindComponent) error {
	obj["keybind"] = c.Keybind
	return s.encodeComponent(obj, c)
}

func (s *Serializer) encodeScore(obj map[string]interface{}, c *ScoreComponent) error {
	score := map[string]interface{}{
		"name":      c.Score.Name,
		"objective": c.Score.Objective,
	}
	if c.Score.Value != "" {
		score["value"] = c.Score.Value
	}
	obj["score"] = score
	return s.encodeComponent(obj, c)
}

func (s *Serializer) encodeSelector(obj map[string]interface{}, c *SelectorComponent) error {
	obj["selector"] = c.Selector
	return s.encodeComponent(obj, c)
}

func (s *Serializer) encodeComponent(obj map[string]interface{}, c Component) error {
	if c.IsBold() {
		obj["bold"] = true
	}

	if c.IsItalic() {
		obj["italic"] = true
	}

	if c.IsUnderlined() {
		obj["underlined"] = true
	}

	if c.IsStrikethrough() {
		obj["strikethrough"] = true
	}

	if c.IsObfuscated() {
		obj["obfuscated"] = true
	}

	if color := c.GetColor(); color != nil {
		if s.ForceLegacyColors {
			if color.Name != "" {
				obj["color"] = color.Name
			} else {
				obj["color"] = FindNearest(*color).Name
			}
		} else {
			obj["color"] = color.String()
		}
	}

	if insertion := c.GetInsertion(); insertion != "" {
		obj["insertion"] = insertion
	}

	if clickEvent := c.GetClickEvent(); clickEvent != nil {
		obj["clickEvent"] = map[string]interface{}{
			"action": clickEvent.Action,
			"value":  clickEvent.Value,
		}
	}

	if hoverEvent := c.GetHoverEvent(); hoverEvent != nil {
		var event = map[string]interface{}{
			"action": hoverEvent.Action,
		}
		if s.LegacyHoverEvent {
			event["value"] = hoverEvent.Contents
		} else {
			event["contents"] = hoverEvent.Contents
		}
		obj["hoverEvent"] = event
	}

	var extra []map[string]interface{}
	for _, e := range c.GetExtra() {
		var comp = make(map[string]interface{})
		if err := s.encode(comp, e); err != nil {
			return err
		}
		extra = append(extra, comp)
	}
	if len(extra) > 0 {
		obj["extra"] = extra
	}

	return nil
}

func (s *Serializer) decode(obj interface{}) (Component, error) {
	switch c := obj.(type) {
	case string:
		return &TextComponent{
			Text: c,
		}, nil
	case map[string]interface{}:
		if _, ok := c["text"]; ok {
			return s.decodeText(c)
		} else if _, ok := c["translate"]; ok {
			return s.decodeTranslatable(c)
		} else if _, ok := c["keybind"]; ok {
			return s.decodeKeybind(c)
		} else if _, ok := c["score"]; ok {
			return s.decodeScore(c)
		} else if _, ok := c["selector"]; ok {
			return s.decodeSelector(c)
		} else {
			return nil, errors.New("json input unmarshalled to unsupported component")
		}
	default:
		return nil, fmt.Errorf("json input unmarshalled to unsupported type %T", obj)
	}
}

func (s *Serializer) decodeText(obj map[string]interface{}) (Component, error) {
	text, ok := obj["text"].(string)
	if !ok {
		return nil, errors.New("text key must be a string")
	}
	return s.decodeComponent(obj, &TextComponent{
		Text: text,
	})
}

func (s *Serializer) decodeTranslatable(obj map[string]interface{}) (Component, error) {
	translate, ok := obj["translate"].(string)
	if !ok {
		return nil, errors.New("translate key must be a string")
	}

	var with []Component
	if array, ok := obj["with"]; ok {
		array, ok := array.([]map[string]interface{})
		if !ok {
			return nil, errors.New("with key must be a array")
		}
		for _, obj := range array {
			c, err := s.decode(obj)
			if err != nil {
				return nil, err
			}
			with = append(with, c)
		}
	}

	return s.decodeComponent(obj, &TranslatableComponent{
		Translate: translate,
		With:      with,
	})
}

func (s *Serializer) decodeKeybind(obj map[string]interface{}) (Component, error) {
	keybind, ok := obj["keybind"].(string)
	if !ok {
		return nil, errors.New("keybind key must be a string")
	}
	return s.decodeComponent(obj, &KeybindComponent{
		Keybind: keybind,
	})
}

func (s *Serializer) decodeScore(obj map[string]interface{}) (Component, error) {
	score, ok := obj["score"].(map[string]interface{})
	if !ok {
		return nil, errors.New("score key must be a object")
	}

	var test Score
	if name, ok := score["name"]; ok {
		name, ok := name.(string)
		if !ok {
			return nil, errors.New("name key must be a string")
		}
		test.Name = name
	}

	if objective, ok := score["objective"]; ok {
		objective, ok := objective.(string)
		if !ok {
			return nil, errors.New("objective key must be a string")
		}
		test.Objective = objective
	}

	if value, ok := score["value"]; ok {
		value, ok := value.(string)
		if !ok {
			return nil, errors.New("value key must be a string")
		}
		test.Value = value
	}

	return s.decodeComponent(obj, &ScoreComponent{
		Score: test,
	})
}

func (s *Serializer) decodeSelector(obj map[string]interface{}) (Component, error) {
	selector, ok := obj["selector"].(string)
	if !ok {
		return nil, errors.New("selector key must be a string")
	}
	return s.decodeComponent(obj, &SelectorComponent{
		Selector: selector,
	})
}

func (s *Serializer) decodeComponent(obj map[string]interface{}, c Component) (Component, error) {
	if bold, ok := obj["bold"]; ok {
		bold, ok := bold.(bool)
		if !ok {
			return nil, errors.New("bold key must be a bool")
		}
		c.SetBold(bold)
	}

	if italic, ok := obj["italic"]; ok {
		italic, ok := italic.(bool)
		if !ok {
			return nil, errors.New("italic key must be a bool")
		}
		c.SetItalic(italic)
	}

	if underlined, ok := obj["underlined"]; ok {
		underlined, ok := underlined.(bool)
		if !ok {
			return nil, errors.New("underlined key must be a bool")
		}
		c.SetUnderlined(underlined)
	}

	if strikethrough, ok := obj["strikethrough"]; ok {
		strikethrough, ok := strikethrough.(bool)
		if !ok {
			return nil, errors.New("strikethrough key must be a bool")
		}
		c.SetStrikethrough(strikethrough)
	}

	if obfuscated, ok := obj["obfuscated"]; ok {
		obfuscated, ok := obfuscated.(bool)
		if !ok {
			return nil, errors.New("obfuscated key must be a bool")
		}
		c.SetObfuscated(obfuscated)
	}

	if color, ok := obj["color"]; ok {
		value, ok := color.(string)
		if !ok {
			return nil, errors.New("color key must be a string")
		}

		color := Color{}
		if strings.HasPrefix(value, "#") {
			color.Hex = strings.TrimPrefix(value, "#")
			if s.ForceLegacyColors {
				color = FindNearest(color)
			}
		} else {
			color = FindByName(value)
		}

		c.SetColor(&color)
	}

	if insertion, ok := obj["insertion"]; ok {
		insertion, ok := insertion.(string)
		if !ok {
			return nil, errors.New("insertion key must be a string")
		}
		c.SetInsertion(insertion)
	}

	if clickEvent, ok := obj["clickEvent"]; ok {
		clickEvent, ok := clickEvent.(map[string]interface{})
		if !ok {
			return nil, errors.New("clickEvent key must be a object")
		}

		event := &ClickEvent{}
		if action, ok := clickEvent["action"]; ok {
			action, ok := action.(string)
			if !ok {
				return nil, errors.New("clickEvent action key must be a string")
			}
			event.Action = ClickAction(action)
		}

		if value, ok := clickEvent["value"]; ok {
			value, ok := value.(string)
			if !ok {
				return nil, errors.New("clickEvent value key must be a string")
			}
			event.Value = value
		}

		c.SetClickEvent(event)
	}

	if hoverEvent, ok := obj["hoverEvent"]; ok {
		hoverEvent, ok := hoverEvent.(map[string]interface{})
		if !ok {
			return nil, errors.New("hoverEvent key must be a object")
		}

		event := &HoverEvent{}
		if action, ok := hoverEvent["action"]; ok {
			action, ok := action.(string)
			if !ok {
				return nil, errors.New("hoverEvent action key must be a string")
			}
			event.Action = HoverAction(action)
		}

		if s.LegacyHoverEvent {
			if value, ok := hoverEvent["value"]; ok {
				event.Contents = value
			}
		} else {
			if contents, ok := hoverEvent["contents"]; ok {
				event.Contents = contents
			}
		}

		c.SetHoverEvent(event)
	}

	if extra, ok := obj["extra"]; ok {
		extra, ok := extra.([]map[string]interface{})
		if !ok {
			return nil, errors.New("extra key must be a array")
		}

		var extras []Component
		for _, obj := range extra {
			c, err := s.decode(obj)
			if err != nil {
				return nil, err
			}
			extras = append(extras, c)
		}

		c.SetExtra(extras)
	}

	return c, nil
}
