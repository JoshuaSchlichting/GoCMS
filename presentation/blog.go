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
