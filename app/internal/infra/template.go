package infra

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	customtemplate "html/template"
	stdtemplate "html/template"

	"upspin.io/errors"
	// customtemplate "github.com/alecthomas/template"
	// blackfriday "gopkg.in/russross/blackfriday.v2"
)

type Template struct {
	templates *stdtemplate.Template
	custom    *stdtemplate.Template
	funcMap   stdtemplate.FuncMap
}

func NewTemplate() *Template {
	funcMap := customtemplate.FuncMap{
		"last":      lastInArray,
		"toJSON":    objectToJSON,
		"asHTML":    stringToHTML,
		"params":    generateDict,
		"unixToUTC": timestampToUTC,
		"b2s":       b2sSI,
		"u2d":       timestampToCustomDate,
		"inc":       Increment,
		"dec":       Decrement,
	}

	templatePagePath := []string{
		"views/front/*.tpl",
		"views/back/*.tpl",
	}
	pagesPath := []string{
		"views/front/pages/*.tpl",
		"views/back/pages/*.tpl",
	}
	var (
		templates = customtemplate.New("template")
		tplPages  = customtemplate.New("page")
	)

	for _, pathGlob := range templatePagePath {
		templates = customtemplate.Must(templates.New("template").Funcs(funcMap).ParseGlob(pathGlob))
	}
	for _, pathGlob := range pagesPath {
		tplPages = customtemplate.Must(tplPages.New("page").Funcs(funcMap).ParseGlob(pathGlob))
	}

	return &Template{
		templates: templates,
		custom:    tplPages,
		funcMap:   funcMap,
	}
}

func (t *Template) JSEscapeString(s string) string {
	return customtemplate.JSEscapeString(s)
}

func (t *Template) Render(w http.ResponseWriter, r *http.Request, status int, name string, data map[string]any) error {
	const op errors.Op = "infra.Render"
	w.WriteHeader(status)

	// add data to main template
	if data == nil {
		data = map[string]any{}
	}
	//data["is_index"] = false
	//data["is_thread"] = false
	data["errors"] = S.Errors
	buffer := bytes.NewBufferString("")
	t.custom.ExecuteTemplate(buffer, name, data)

	// maybe I shouldn't hardcode template names in is_index and is_thread
	content := map[string]any{
		"navigation": boards.BoardList(S.Conn),
		"page":       buffer.String(),
		"is_index":   strings.HasSuffix(name, "thread_list"),
		"is_thread":  strings.HasSuffix(name, "post_list"),
		"data":       S.Data,
	}

	baseTplName := "template"
	environment := "front"
	if strings.HasPrefix(name, "back") {
		environment = "back"
	}
	baseTplName = strings.Join([]string{environment, baseTplName}, "/")

	err := t.templates.ExecuteTemplate(w, baseTplName, content)
	if err != nil {
		log.Fatalf("Could not render %s (%s)", name, err)
	}

	return err
}

func (t *Template) StringToHTML(s string) stdtemplate.HTML {
	return stdtemplate.HTML(s)
}

// func (t *Template) MarkdownToHTML(s string) stdtemplate.HTML {
// 	return stdtemplate.HTML(blackfriday.Run([]byte(s)))
// }
