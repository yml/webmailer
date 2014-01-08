package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

const (
	base          = `{{ define "base" }}<html><body>{{ template "content" }}</body></html>{{ end }}`
	emailForm     = `{{ define "content" }}<h1>Send us an email</h1><form action="/send" method="post"><label>body</label><textarea name="body"></textarea><button type="submit">Send</button></form>{{ end }}`
	thanks        = `{{ define "content" }}<h1>Thank you</h1>{{ end }}`
	emailTemplate = `From: {{.From}}
To: {{.To}}}
Subject: {{.Subject}}

{{.Body}}

Sincerely,

{{.From}}`
)

var (
	username, pwd, host, port, addr string
	templates                       = make(map[string]*template.Template)
)

func init() {
	username = os.Getenv("WEBMAILER_SMTP_USERNAME")
	pwd = os.Getenv("WEBMAILER_SMTP_PWD")
	host = os.Getenv("WEBMAILER_SMTP_HOST")
	port = os.Getenv("WEBMAILER_SMTP_PORT")
	addr = os.Getenv("WEBMAILER_HTTP_ADDR")
}

func newTmpl(name string, str string) *template.Template {
	tmpl := template.Must(template.New(name).Parse(base))
	tmpl = template.Must(tmpl.Parse(str))
	return tmpl
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates["emailForm"].ExecuteTemplate(w, "base", "content")
}

func thanksHandler(w http.ResponseWriter, r *http.Request) {
	templates["thanks"].ExecuteTemplate(w, "base", "content")
}

func sendMail(username, pwd, host, port, body string) (err error) {
	var msg bytes.Buffer
	ctx := make(map[string]string)
	ctx["From"] = username
	ctx["To"] = username
	ctx["Subject"] = "Mail from Go app"
	ctx["Body"] = body

	tmpl := template.Must(template.New("email").Parse(emailTemplate))
	err = tmpl.Execute(&msg, ctx)
	if err != nil {
		return
	}
	auth := smtp.PlainAuth("", username, pwd, host)
	err = smtp.SendMail(host+":"+port, auth, username,
		[]string{username}, msg.Bytes())

	return
}

func sendMailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body := r.PostFormValue("body")
		err := sendMail(username, pwd, host, port, body)
		if err != nil {
			log.Println("[Error] An error occured while sending the email:", err)
		}
		http.Redirect(w, r, "/thanks", http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func main() {
	templates["emailForm"] = newTmpl("emailForm", emailForm)
	templates["thanks"] = newTmpl("thanks", thanks)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/send", sendMailHandler)
	http.HandleFunc("/thanks", thanksHandler)
	log.Println("Starting an http server listening on :", addr)
	http.ListenAndServe(addr, nil)
}
