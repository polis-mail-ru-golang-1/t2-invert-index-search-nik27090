package handlers

import (
	"html/template"
	"net/http"

	"../invertIndex"
	"go.uber.org/zap"
)

func AddPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			zap.S().Info("Failed addPage")
		}
		/*fileName := r.FormValue("filename")
		content := r.FormValue("content")*/
	}
	tmpl := template.Must(template.ParseFiles("templates/addFile.html"))
	tmpl.Execute(w, nil)
}

func SearchPage(w http.ResponseWriter, r *http.Request) {

	q := r.FormValue("q")
	zap.S().Info("Search phrase: ", q)
	if q != "" {
		endMap := invertIndex.Find(invertIndex.InIn, q, invertIndex.SliceFiles)
		endSlice := invertIndex.SortSearch(endMap, w)
		tmpl := template.Must(template.ParseFiles("templates/searchPage.html", "templates/dynamicList.html"))
		tmpl.Execute(w, endSlice)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/mainPage.html"))
	tmpl.Execute(w, nil)
}
