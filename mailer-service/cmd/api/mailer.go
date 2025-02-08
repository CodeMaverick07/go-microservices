package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain string
	Host string
	Port int
	Username string
	Password string
	Encryption string
	FromAddress string
	FromName string
}

type Message struct {
	From      	string
	FromName  	string
	To        	string
	Subject   	string
	Attachments []string
	Data        any 
	DataMap     map[string]any
}

func (m *Mail) sendSMTPMessage(msg Message) (error) {
	if (msg.From == "") {
		msg.From = m.FromAddress
	}
	if (msg.FromName == "") {
		msg.FromName = m.FromName
	}
	data := map[string]interface{}{
		"message": msg.Data,
	}
	msg.DataMap = data
    formmatedMessage,err := m.buildHTMLMessage(msg)
	if err != nil {
		
		return err 
	}
	plainMessage,err := m.buildPlainTextMessage(msg)
	if err != nil {
		
		return err 
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.KeepAlive = false
	 
	smtpClient,err := server.Connect()
	if err != nil {
		
		return err 
	}

	email:= mail.NewMSG()
	email.SetFrom(msg.From).
		  AddTo(msg.To).
		  SetSubject(msg.Subject).
		  SetBody(mail.TextPlain,plainMessage).
		  AddAlternative(mail.TextHTML,formmatedMessage)

    if len(msg.Attachments) > 0 {
		for _,attachment := range msg.Attachments {
			email.AddAttachment(attachment)
		}
	}
	err = email.Send(smtpClient)
	if err != nil {
		
		return err
	}
	return nil
}
func (m *Mail) buildHTMLMessage(msg Message) (string,error) {
     templateToRender := "./templates/mail.html.gohtml"
	 t,err := template.New("email-html").ParseFiles(templateToRender)
	 if err != nil {
		
		 return "",err
	 }
	 var tpl bytes.Buffer
	 if err = t.ExecuteTemplate(&tpl,"body",msg.DataMap); err != nil {
		
		 return "",err
	 }
	 formattedMessage := tpl.String()
	 formattedMessage,err = m.inlineCSS(formattedMessage)

     if err != nil {
		
		 return "",err
	 }
	 return formattedMessage,nil
}
func (m *Mail) buildPlainTextMessage(msg Message) (string,error) {
     templateToRender := "./templates/mail.plain.gohtml"
	 t,err := template.New("email-plain").ParseFiles(templateToRender)
	 if err != nil {
		
		 return "",err
	 }
	 var tpl bytes.Buffer
	 if err = t.ExecuteTemplate(&tpl,"body",msg.DataMap); err != nil {
		
		 return "",err
	 }
	 plainMessage:= tpl.String()
	 return plainMessage,nil
}
func (m *Mail) inlineCSS(body string) (string,error) {
	options := premailer.Options{
		RemoveClasses: false,
		CssToAttributes: false,
		KeepBangImportant: true,
	}
	prem,err := premailer.NewPremailerFromString(body,&options)
	if err != nil {
		
		return "",err
	}
	html,err := prem.Transform()
	if err != nil {
		
		return "",err
	}
	return html,nil
}
func(m *Mail) getEncryption(s string) mail.Encryption { 
	switch s {
	case "tls":
		return mail.EncryptionTLS
	case "ssl":
		return mail.EncryptionSSL
	case "none","":
		return mail.EncryptionNone 
	default:
		return mail.EncryptionSTARTTLS 
	}
}


