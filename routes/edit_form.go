package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type FormField struct {
	Name  string
	Type  string
	Value string
}

func GenerateForm(title string, fields []FormField, hxMethod, hxURL, idPrefix string) template.HTML {
	tmpl := `
		<div class="card">
			<div class="card-body">
				<h5 class="card-title">{{.Title}}</h5>
				<form id="{{.FormID}}" class="row g-3" hx-{{.HxMethod}}="{{.HxURL}}">
					{{range $field := .Fields}}
						<div class="col-12">
							<label class="form-label">{{$field.Name}}</label>
							<input type="{{$field.Type}}" class="form-control" id="{{$.IDPrefix}}{{$field.Name}}" name="{{$field.Name}}" value="{{$field.Value}}" {{if eq (ToLower $field.Name) "id"}}readonly{{end}}>
						</div>
					{{end}}
					<div class="text-center">
						<button type="submit" class="btn btn-primary">Submit</button>
						<button type="reset" class="btn btn-secondary">Reset</button>
					</div>
				</form>
			</div>
		</div>
		<script type="text/javascript">

		</script>`

	// Create a map of parameters to pass to the template
	data := struct {
		Title    string
		Fields   []FormField
		FormID   string
		IDPrefix string
		HxMethod template.JS
		HxURL    string
	}{
		Title:    title,
		Fields:   fields,
		FormID:   idPrefix + "_form",
		IDPrefix: idPrefix,
		HxMethod: template.JS(hxMethod),
		HxURL:    hxURL,
	}

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	// Execute the template
	var tpl bytes.Buffer
	formTmpl := template.Must(template.New("").Funcs(funcMap).Parse(tmpl))
	err := formTmpl.Execute(&tpl, data)
	if err != nil {
		fmt.Println("Error generating form:", err)
		return ""
	}

	return template.HTML(tpl.String())
}
