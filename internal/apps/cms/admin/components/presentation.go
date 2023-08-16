// Package presentation provides the presentation layer for the gocms application.
// The Presentor struct is used to render the templates, writing them directly to
// the http.ResponseWriter passed to NewPresentor.
package components

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"
)

type ClickableTable struct {
	TableID      template.JS
	Table        []map[string]interface{}
	CallbackFunc template.JS
	JavaScript   template.JS
}

type FormField struct {
	Name  string
	Type  string
	Value string
}

type Presentor struct {
	template *template.Template
	writer   io.Writer
}

// NewPresentor returns a new Presentor struct who's methods will write directly
// to the http.ResponseWriter passed to this constructor.
func NewPresentor(t *template.Template, w io.Writer) *Presentor {
	return &Presentor{
		template: t,
		writer:   w,
	}
}

func (p *Presentor) CreateItemFormHTML(formID, formTitle, apiEndpoint, refreshURL string, formFields []FormField) error {
	return p.template.ExecuteTemplate(p.writer, "modify_item_form", map[string]interface{}{
		"Form": generateForm(
			formTitle,
			formFields,
			"post",
			apiEndpoint,
			formID,
			getSubmitResetButtonDiv(),
		),
		"EventTargetID": formID + "_form",
		"FormID":        formID,
		"RefreshURL":    template.JS(refreshURL),
	})
}

func (p *Presentor) EditListItemHTML(formID, formTitle, apiEndpoint, apiCallType, refreshURL string, formFields []FormField, dataMap []map[string]interface{}) error {

	jsSetFormElements := ""
	for _, field := range formFields {
		jsSetFormElements += fmt.Sprintf(`
			// loop over form fields and add them
			document.getElementById("%[1]s%s").value = formData["%[2]s"];
		`, formID, field.Name)
	}
	jsSetApiUrl := fmt.Sprintf(`
		document.getElementById("%[1]s_form").setAttribute("hx-%[3]s", "%[2]s/" + formData["ID"]);
		htmx.process(document.getElementById("%[1]s_form"));`, formID, apiEndpoint, apiCallType)
	tableID := formID + "_table"

	initDataTableCode := "" //fmt.Sprintf(`let table = new DataTable('#%s', {});`, tableID)

	err := p.template.ExecuteTemplate(p.writer, "modify_item_form", map[string]interface{}{
		"Form": generateForm(
			formTitle,
			formFields,
			"put",
			apiEndpoint,
			formID,
			getSubmitResetButtonDiv(),
		),
		"EventTargetID": formID + "_form",
		"FormID":        formID,
		"RefreshURL":    template.JS(refreshURL),
		"ClickableTable": &ClickableTable{
			TableID:      template.JS(tableID),
			Table:        dataMap,
			CallbackFunc: template.JS("setItemInForm"),
			JavaScript:   template.JS(fmt.Sprintf(string(setItemInFormJS), jsSetFormElements+jsSetApiUrl, tableID)) + template.JS(initDataTableCode),
		},
	})
	return err
}

func (p *Presentor) DeleteItemFormHTML(formID, formTitle, apiEndpoint, refreshURL, setItemAdditionalJS string, formFields []FormField, dataMap []map[string]interface{}) error {
	tableID := formID + "_table"
	jsSetFormElements := ""
	for _, field := range formFields {
		jsSetFormElements += fmt.Sprintf(`
			// loop over form fields and add them
			document.getElementById("%[1]s%s").value = formData["%[2]s"];
		`, formID, field.Name)
	}

	err := p.template.ExecuteTemplate(p.writer, "modify_item_form", map[string]interface{}{
		"Form": generateForm(
			formTitle,
			formFields,
			"delete",
			apiEndpoint,
			formID,
			`<button id="deleteItemButton" hx-confirm="Are you sure you want to delete this object?" type="button" class="btn btn-danger">Delete</button>`,
		),
		"FormID":        formID,
		"RefreshURL":    template.JS(refreshURL),
		"EventTargetID": "deleteItemButton",
		"ClickableTable": &ClickableTable{
			TableID:      template.JS(tableID),
			Table:        dataMap,
			CallbackFunc: template.JS("setItemInForm"),
			JavaScript:   template.JS(fmt.Sprintf(string(setItemInFormJS), jsSetFormElements+setItemAdditionalJS, tableID)),
		},
	})
	return err
}

func (p *Presentor) CreateBlogHTML() error {
	return p.template.ExecuteTemplate(p.writer, "create_blog", map[string]interface{}{
		"CreateBlogForm": generateForm(
			"Create Blog",
			[]FormField{
				{
					Name:  "title",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "description",
					Type:  "text",
					Value: "",
				},
			},
			"post",
			"/api/v1/blogs",
			"create-blog-form",
			getSubmitResetButtonDiv(),
		),
		"FormID": "create-blog-form",
		"ImageURLs": []string{
			"https://picsum.photos/200/300?random=1",
			"https://picsum.photos/200/300?random=2",
			"https://picsum.photos/200/300?random=3",
			"https://picsum.photos/200/300?random=4",
			"https://picsum.photos/200/300?random=5",
			"https://picsum.photos/200/300?random=6",
			"https://picsum.photos/200/300?random=7",
			"https://picsum.photos/200/300?random=8",
			"https://picsum.photos/200/300?random=9",
			"https://picsum.photos/200/300?random=10",
		},
	})
}

func getSubmitResetButtonDiv() string {
	return `
		<div class="text-center">
		<button type="submit" class="btn btn-primary">Submit</button>
		<button type="reset" class="btn btn-secondary">Reset</button>
		</div>`
}

func generateForm(title string, fields []FormField, hxMethod, hxURL, idPrefix, buttonDiv string) template.HTML {
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
					{{.ButtonDiv}}
				</form>
			</div>
		</div>`

	// Create a map of parameters to pass to the template
	data := struct {
		Title     string
		Fields    []FormField
		FormID    string
		IDPrefix  string
		HxMethod  template.JS
		HxURL     string
		ButtonDiv template.HTML
	}{
		Title:     title,
		Fields:    fields,
		FormID:    idPrefix + "_form",
		IDPrefix:  idPrefix,
		HxMethod:  template.JS(hxMethod),
		HxURL:     hxURL,
		ButtonDiv: template.HTML(buttonDiv),
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

const setItemInFormJS template.JS = (`function setItemInForm() {
	console.log("setItemInForm->");
	let formData = getRowData("%[2]s", "ID", %[2]sSelectedRow);
	console.log(formData);
	%[1]s
}`)
