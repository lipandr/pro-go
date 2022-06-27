package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Rsvp struct {
	Name, Email, Phone string
	WillAttend         bool
}

type formData struct {
	*Rsvp
	Errors []string
}

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("layout.html", name+".html")
		if err != nil {
			panic(err)
		}
		templates[name] = t
		fmt.Println("loaded template ", index, name)
	}
}

func welcomeHandler(w http.ResponseWriter, _ *http.Request) {
	err := templates["welcome"].Execute(w, nil)
	if err != nil {
		return
	}
}

func listHandler(w http.ResponseWriter, _ *http.Request) {
	_ = templates["list"].Execute(w, responses)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = templates["form"].Execute(w, formData{
			&Rsvp{}, []string{},
		})
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			return
		}
		responseData := Rsvp{
			Name:       r.Form["name"][0],
			Email:      r.Form["email"][0],
			Phone:      r.Form["phone"][0],
			WillAttend: r.Form["willattend"][0] == "true",
		}
		var errors []string
		if responseData.Name == "" {
			errors = append(errors, "Please enter your name")
		}
		if responseData.Email == "" {
			errors = append(errors, "Please enter your email address")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Please enter your phone number")
		}
		if len(errors) > 0 {
			_ = templates["form"].Execute(w, formData{
				Rsvp: &responseData, Errors: errors,
			})
		} else {
			responses = append(responses, &responseData)

			if responseData.WillAttend {
				_ = templates["thanks"].Execute(w, responseData.Name)
			} else {
				_ = templates["sorry"].Execute(w, responseData.Name)
			}
		}
	}
}

func main() {
	loadTemplates()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
