package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/syamimhazmi/snippetbox/internal/models"
	"github.com/syamimhazmi/snippetbox/internal/validator"
)

type SignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type LoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type SnippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *Application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = SignupForm{}

	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *Application) signupPost(w http.ResponseWriter, r *http.Request) {
	var form SignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(
		validator.NotBlank(form.Name),
		"name",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.NotBlank(form.Email),
		"email",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.Matches(form.Email, validator.EmailRX),
		"email",
		"This field must be a valid email address",
	)
	form.CheckField(
		validator.NotBlank(form.Password),
		"password",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.MinChars(form.Password, 8),
		"password",
		"This field must be at least 8 charactes long",
	)

	if !form.Valid() {
		data := app.newTemplateData(r)

		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)

		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(
		r.Context(),
		"flash",
		"Your signup was successful. Please log in.",
	)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = LoginForm{}

	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *Application) loginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *Application) logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = SnippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *Application) snippetStore(w http.ResponseWriter, r *http.Request) {
	var form SnippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(
		validator.NotBlank(form.Title),
		"title",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.MaxChars(form.Title, 100),
		"title",
		"This field cannot be more than 100 characters long",
	)
	form.CheckField(
		validator.NotBlank(form.Content),
		"content",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.PermittedValue(form.Expires, 1, 7, 365),
		"expires",
		"This field must equal 1, 7 or 365",
	)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
}
