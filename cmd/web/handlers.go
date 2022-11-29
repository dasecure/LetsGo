package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.dasecure.com/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// log.Println("im at home")
	// panic("oops! something went wrong!")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serveError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

// for _, snippet := range snippets {
// 	fmt.Fprintf(w, "%+v\n", snippet)
// }

// files := []string{
// 	"./ui/html/partials/nav.tmpl",
// 	"./ui/html/pages/home.tmpl",
// 	"./ui/html/base.tmpl",
// }

// ts, err := template.ParseFiles(files...)
// if err != nil {
// 	log.Println(err.Error())
// 	// http.Error(w, "Internal server error", 500)
// 	app.serveError(w, err)
// 	return
// }

// data := &templateData{
// 	Snippets: snippets,
// }

// err = ts.ExecuteTemplate(w, "base", data)
// if err != nil {
// 	log.Println(err.Error())
// 	// http.Error(w, "Internal server error", 500)
// 	app.serveError(w, err)
// }
// app.render(w, http.StatusOK, "home.tmpl", &templateData{
// 	Snippets: snippets,
// })
// }

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// log.Println("I entered view request")
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}

	// log.Println(id)
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serveError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

// app.render(w, http.StatusOK, "view.tmpl", &templateData{
// 	Snippet: snippet,
// })

// fmt.Fprintf(w, "%+v", snippet)

// files := []string{
// 	"./ui/html/partials/nav.tmpl",
// 	"./ui/html/pages/view.tmpl",
// 	"./ui/html/base.tmpl",
// }

// ts, err := template.ParseFiles(files...)
// if err != nil {
// 	log.Println(err.Error())
// 	// http.Error(w, "Internal server error", 500)
// 	app.serveError(w, err)
// 	return
// }

// data := &templateData{
// 	Snippet: snippet,
// }

// err = ts.ExecuteTemplate(w, "base", data)
// if err != nil {
// 	log.Println(err.Error())
// 	// http.Error(w, "Internal server error", 500)
// 	app.serveError(w, err)
// }

// func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		w.Header().Set("Allow", http.MethodPost)
// 		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		app.clientError(w, http.StatusMethodNotAllowed)
// 		return
// 	}
// 	title := "Slow snail"
// 	content := "The slow snail was bloody irritating\n, it took so long to get to the finish line.\n\nJack Neo"
// 	expires := 7
// 	id, err := app.snippets.Insert(title, content, expires)
// 	if err != nil {
// 		app.serveError(w, err)
// 		return
// 	}
// 	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
// 	// w.Write([]byte("Create a new snippet"))
// }

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet"))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "Create a new snippet"
	content := "The slow snail was bloody irritating\n, it took"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serveError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
