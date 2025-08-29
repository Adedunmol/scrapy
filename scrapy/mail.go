package scrapy

import (
	"bytes"
	"fmt"
	"github.com/Adedunmol/scrapy/core"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

const Subject = "These are the job postings for today."
const Template = "jobs"

func SendMail(email string, parsedJobs []*core.Job) error {

	from := strings.TrimSpace(os.Getenv("FROM_EMAIL"))
	password := strings.TrimSpace(os.Getenv("FROM_EMAIL_PASSWORD"))
	smtpAddr := strings.TrimSpace(os.Getenv("SMTP_ADDR"))
	smtpPort := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	adminEmail := strings.TrimSpace(os.Getenv("ADMIN_EMAIL"))

	auth := smtp.PlainAuth("", from, password, smtpAddr)

	// Build full headers
	headers := make(map[string]string)
	headers["From"] = adminEmail
	headers["To"] = strings.Join([]string{email}, ", ")
	headers["Subject"] = Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	htmlBody, err := parseTemplate(parsedJobs)
	if err != nil {
		err = fmt.Errorf("error parsing template: %v", err)
		fmt.Println(err)
		return err
	}

	// Construct email message
	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")   // Blank line between headers and body
	msg.WriteString(htmlBody) // HTML content

	// Send email
	return smtp.SendMail(
		smtpAddr+":"+smtpPort,
		auth,
		from,
		[]string{email},
		[]byte(msg.String()),
	)
}

func parseTemplate(data []*core.Job) (string, error) {

	currDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}
	templatePath := filepath.Join(currDir, Template+".html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return rendered.String(), nil
}
