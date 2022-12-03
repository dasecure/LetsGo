package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.dasecure.com/internal/models"
	"snippetbox.dasecure.com/internal/validator"
)

type userSignForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// user authentication
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
	// fmt.Fprintln(w, "Display a signup form")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must exceed 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serveError(w, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	// fmt.Fprintln(w, "Create new user")
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
	// fmt.Fprintln(w, "Display a login form")
}

// fmt.Fprintln(w, "Display a login form")

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// log.Println("before authentication")
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or Password is incorrect.")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serveError(w, err)
		}
		return
	}
	log.Printf("After authentication.. Handlers.go %d", id)
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serveError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	// fmt.Fprintln(w, "Auth and login user")
}
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serveError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "AuthenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You are now logged out")

	http.Redirect(w, r, "/", http.StatusSeeOther)
	// fmt.Fprintln(w, "Logout the user")
}

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

	flash := app.sessionManager.PopString(r.Context(), "flash")
	data := app.newTemplateData(r)

	data.Snippet = snippet
	data.Flash = flash

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

type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	// FieldErrors map[string]string
	validator.Validator `form:"-"`
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	log.Print("Enter snippet created")
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// title := "Create a new snippet"
	// content := "The slow snail was bloody irritating\n, it took"
	// expires := 7
	// err := r.ParseForm()
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")
	// expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// form := snippetCreateForm{
	// 	Title:   r.PostForm.Get("title"),
	// 	Content: r.PostForm.Get("content"),
	// 	Expires: expires,
	// 	// FieldErrors: map[string]string{},
	// }

	// fieldErrors := make(map[string]string)
	form.CheckField(validator.NotBlank(form.Title), "title", "Title is required.")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "Title must not exceed 100 characters.")
	// if strings.TrimSpace(form.Title) == "" {
	// 	form.FieldErrors["title"] = "This field cannot be empty"
	// } else if utf8.RuneCountInString(form.Title) > 100 {
	// 	form.FieldErrors["title"] = "Title cannot be longer than 100 characters"
	// }
	form.CheckField(validator.NotBlank(form.Content), "content", "Content cant be empty")
	// if strings.TrimSpace(form.Content) == "" {
	// 	form.FieldErrors["content"] = "This field cannot be empty"
	// }
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must be 1, 7 or 365")
	// if expires != 1 && expires != 7 && expires != 365 {
	// 	form.FieldErrors["expires"] = "This field must be 1, 7 or 365"
	// }

	if !form.Valid() {
		log.Println("Theres an error")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serveError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
