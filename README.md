# Webmailer

This is an app that send emails, it exposes 3 urls:

* **/** index
* **/send** send email that is POSTed
* **/thanks** thank you page


## Install

Download the webmailer code from github using `go get`

```
go get -u github.com/yml/webmailer
```

In a terminal set the environ variables required for the web mailer

```
export WEBMAILER_SMTP_USERNAME="username@gmail.com"
export WEBMAILER_SMTP_PWD="secret"
export WEBMAILER_SMTP_HOST="smtp.gmail.com"
export WEBMAILER_SMTP_PORT="587"
export WEBMAILER_HTTP_ADDR=":8080"
```

Start `webmailer`
