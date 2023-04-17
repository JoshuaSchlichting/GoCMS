package presentation

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
	})
}
