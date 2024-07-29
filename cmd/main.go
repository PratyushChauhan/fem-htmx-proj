package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func newtemplate() *Templates {
	return &Templates{templates: template.Must(template.ParseGlob("views/*.html"))}
}

type Count struct {
	Count int
}
type Data struct {
	Contacts []Contact
}

func (d *Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

type Contact struct {
	Name  string
	Email string
}

// Values is a map of form values for the input with an error.
// Errors is a map of errors for the particular field eg : "email" -> "This already exists".
type FormData struct {
	Values map[string]string
	Errors map[string]string
}
type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{Data: newData(), Form: newFormData()}
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}
func newContact(name, email string) Contact {
	return Contact{Name: name, Email: email}
}
func newData() Data {
	return Data{Contacts: []Contact{
		newContact("John Doe", "johndoe@gmail.com"),
		newContact("Clara", "clara@gmail.com"),
		newContact("Jenny", "jenny@gmail.com"),
	}}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = newtemplate()
	page := newPage()
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})
	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		if page.Data.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "This email already exists"
			return c.Render(422, "form", formData)
		}
		page.Data.Contacts = append(page.Data.Contacts, newContact(name, email))
		return c.Render(200, "display", page.Data)
	})
	e.Logger.Fatal(e.Start(":42069"))

}
