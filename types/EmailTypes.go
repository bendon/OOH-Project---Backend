package types

import "html/template"

type EmailPayload struct {
	Name         string        `json:"name"`
	Subject      string        `json:"subject"`
	MailTo       string        `json:"mailTo"`
	TemplateFile string        `json:"templateFile"`
	Body         template.HTML `json:"body"`
	Link         *string       `json:"link"`
	Code         *string       `json:"code"`
	Attatchments []string      `json:"attachments"`
}

type InvoiceData struct {
	Data1       interface{}
	Data2       interface{}
	Data3       interface{}
	Data4       interface{}
	Data5       interface{}
	Data6       interface{}
	Data7       interface{}
	TotalAmount *float64
	Logo1       *string
	Logo2       *string
	Logo3       *string
	Logo4       *string
	Text1       *string
	Text2       *string
	Text3       *string
	Text4       *string
	Text5       *string
}

type SendEmailRequest struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Subject string  `json:"subject"`
	Body    string  `json:"body"`
	Link    *string `json:"link"`
}
