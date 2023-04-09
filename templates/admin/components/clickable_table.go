package components

import "html/template"

type ClickableTable struct {
	TableID      template.JS
	Table        []map[string]interface{}
	CallbackFunc template.JS
	JavaScript   template.JS
}
