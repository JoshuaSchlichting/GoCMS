package components

import "html/template"

type ClickableTable struct {
	TableID      template.JS
	Table        []map[string]interface{}
	PostURL      string
	CallbackFunc template.JS
	JavaScript   template.JS
}
