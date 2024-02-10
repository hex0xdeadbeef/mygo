package d4_texttemplates

import (
	"golang/pkg/chapters/chapter4/d_structs/d3_json/github"
	h "html/template"
	"log"
	"os"
	t "text/template"
	"time"
)

/*
1. A template is a string/file containing one or more portions enclosed in double braces {{...}}
2. {{...}} - action. In each action we can:
	1) Print values
	2) Select struct fields
	3) Call functions and methods
	4) Express control flow (if-else/range loops)
	5) Instantiate other templates
3. Producing output with a template is a two-step process.
	1) We must parse the template into a suitable internal representation.
	2) Execute specific inputs.
4. Template creation methods chain is:
	1) "template.New" creates and returns a template
	2) "Funcs" adds all the provided functions to the set of functions accessible within this template
	3) "Parse" is called on the result
5. Because templates are usually fixed at compile time, failure to parse a template indicates a fatal bug in the
program.
	1) "template.Must" function makes error handling more convinient:
		1) It accepts a template and an error
		2) Checks that the error is nil/panic
		3) Returns the template
6. Created report works in the following way:
	1) Call "Must" to spot the errors
	2) Call "New" to create named object
	3) Call "Parse" to parse a template

7. To output the data with our template we call "object.Execute(destiny, data source)" with preceding if err state
ment
8. While working with HTML templating we should use template.HTML type to suppresse malicious metacharacters escaping
to prevent ourselves from an injection attack. Using this type we say that our data is trusted and doesn't need to be
escaped.
*/

const (
	templ1 = `{{.TotalCount}} issues:
	{{range .Items}}
	-----------------------------------------------------
	Number: {{.Number}}
	User: {{.User.Login}}
	Title: {{.Title | printf "%.65s"}}
	Age: {{.CreatedAt | daysAgo}} days
	{{end}}`
)

var (
	report = t.
		Must(
			t.New("report").
				Funcs(t.FuncMap{"daysAgo": daysAgo}).
				Parse(templ1))

	issueList = h.
			Must(h.
				New("issueList").
				Parse(
				`
						<h1>{{.TotalCount}} issues:</h1>
						<table>
						<tr style = 'text-align: left'>
							<th>#</th>
							<th>State</th>
							<th>User</th>
							<th>Title</th>
						</tr>
						{{range .Items}}
						<tr>
							<td><a href='{{.HTMLURL}}'>{{.Number}}</td>
							<td>{{.State}}</td>
							<td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
							<td><a href ='{{.HTMLURL}}'>{{.Title}}</a></td>
						</tr>
						{{end}}
						</table>
						`))
)

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func TextTemplateUsing() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}

func HMTLTemplateUsing() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("issues1.html", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := issueList.Execute(file, result); err != nil {
		log.Fatal(err)
	}
}

func StringAndHTMLTemplate() {
	const templ = `<p>A: {{.A}}</p><p>B: {{.B}}</p>`
	t := t.Must(t.New("escape").Parse(templ))

	var data struct {
		A string // untrusted plain text
		B h.HTML // trusted HTML
	}
	data.A = "<b>Hello!</b>"
	data.B = "<b>Hello!</b>"

	if err := t.Execute(os.Stdout, data); err != nil {
		log.Fatal(err)
	}
}
