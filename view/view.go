package view

import (
	"html/template"
	"io"

	"../invertIndex"
)

type View struct {
	MainPage   *template.Template
	SearchPage *template.Template
}

type SearchResult struct {
	[]invertIndex.Result
}

func New() (View, error) {
	v := View{}
	var err error

	v.SearchPage, err = template.ParseFiles("templates/searchPage.html", "templates/dynamicList.html")
	if err != nil {
		return v, err
	}

	v.MainPage, err = template.ParseFiles("templates/mainPage.html")
	if err != nil {
		return v, err
	}

	return v, nil
}

func (v View) SearchView(endSlice []invertIndex.Result, w io.Writer) {
	v.SearchPage.ExecuteTemplate(w, "SearchView",
		struct {
			Result []SearchResult
		}{
			Result: endSlice,
		})
}

func (v View) MainView(w io.Writer) {
	v.SearchPage.ExecuteTemplate(w, "MainPage", nil)
}
