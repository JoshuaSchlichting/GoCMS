package presentation

import (
	"fmt"
	"html/template"
	"io"

	"github.com/joshuaschlichting/gocms/templates/components"

	"bytes"
	"strings"
)

type FormField struct {
	Name  string
	Type  string
	Value string
}

type Presentor struct {
	template *template.Template
	writer   io.Writer
}

func NewPresentor(t *template.Template, w io.Writer) *Presentor {
	return &Presentor{
		template: t,
		writer:   w,
	}
}

func (p *Presentor) GetEditListItemHTML(formID, formTitle, apiEndpoint, apiCallType, refreshURL string, formFields []FormField, dataMap []map[string]interface{}) error {

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

	err := p.template.ExecuteTemplate(p.writer, "edit_item_form", map[string]interface{}{
		"EditUserForm": GenerateForm(
			"Edit User Form",
			formFields,
			"put",
			apiEndpoint,
			formID,
		),
		"FormID":     formID,
		"RefreshURL": template.JS(refreshURL),
		"ClickableTable": &components.ClickableTable{
			TableID:      template.JS(tableID),
			Table:        dataMap,
			CallbackFunc: template.JS("setItemInForm"),
			JavaScript: template.JS(fmt.Sprintf(`
				function getRowData(tableId, columnName, columnValue) {
					console.log("getRowData-> params: " + tableId  + columnName + columnValue);
					// Get the table using its ID
					const table = document.getElementById(tableId);
				
					// Get the table headers
					const headers = table.getElementsByTagName("th");
				
					// Get the index of the target column
					let targetColumnIndex;
					for (let i = 0; i < headers.length; i++) {
						if (headers[i].textContent === columnName) {
							targetColumnIndex = i;
							break;
						}
					}
				
					// Get the table rows
					const rows = table.getElementsByTagName("tr");
				
					// Loop through each row
					for (let i = 0; i < rows.length; i++) {
					const cells = rows[i].getElementsByTagName("td");
				
					// Check if the target column exists in the row
					if (cells[targetColumnIndex]) {
						// Check if the value of the target column matches the columnValue
						if (cells[targetColumnIndex].textContent === columnValue) {
						// Get the header names
						const headerNames = Array.from(headers).map(header => header.textContent);
				
						// Get the cell values
						const cellValues = Array.from(cells).map(cell => cell.textContent);
				
						// Combine the header names and cell values into an object
						const rowData = headerNames.reduce((obj, headerName, index) => {
							obj[headerName] = cellValues[index];
							return obj;
						}, {});
						return rowData;
					}}}
					// Return null if the row is not found
					return null;
				}
				function setItemInForm() {
					console.log("setItemInForm->");
					// get row data where row id == %[2]sSelectedRow
					let formData = getRowData("%[2]s", "ID", %[2]sSelectedRow);
					// get the form
					console.log(formData);
					// prefill the form with id edit_user_form
					%[1]s
				}
			`, jsSetFormElements+jsSetApiUrl, tableID)),
		},
	})
	return err
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
