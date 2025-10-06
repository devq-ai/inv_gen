package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"path/filepath"
	"strings"

	"github.com/devqai/invoice-tui/models"
)

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	GmailAddress  string
	GmailPassword string
	SMTPHost      string
	SMTPPort      string
}

// loadEmailConfig loads email config from .env
func loadEmailConfig() (*EmailConfig, error) {
	envPath := filepath.Join("..", ".env")
	data, err := ioutil.ReadFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .env: %v", err)
	}

	config := &EmailConfig{
		SMTPHost: "smtp.gmail.com",
		SMTPPort: "587",
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "GMAIL_ADDRESS=") {
			config.GmailAddress = strings.TrimSpace(strings.TrimPrefix(line, "GMAIL_ADDRESS="))
		}
		if strings.HasPrefix(line, "GMAIL_APP_PASSWORD=") {
			config.GmailPassword = strings.TrimSpace(strings.TrimPrefix(line, "GMAIL_APP_PASSWORD="))
		}
	}

	if config.GmailAddress == "" || config.GmailPassword == "" {
		return nil, fmt.Errorf("missing email credentials in .env")
	}

	return config, nil
}

// sendInvoiceEmail sends invoice via SMTP with PDF attachment
func sendInvoiceEmail(invoice *models.Invoice, pdfPath string) error {
	config, err := loadEmailConfig()
	if err != nil {
		return err
	}

	// Read PDF file
	pdfData, err := ioutil.ReadFile(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to read PDF: %v", err)
	}

	// Build email
	to := []string{"infoobjects@bill.com", "timesheets@infoobjects.com", config.GmailAddress}
	subject := fmt.Sprintf("Weekly Invoice %s - Dion Edge", invoice.InvoiceNumber)

	boundary := "----=_Part_0_123456789.123456789"
	
	headers := make(map[string]string)
	headers["From"] = config.GmailAddress
	headers["To"] = "infoobjects@bill.com, timesheets@infoobjects.com"
	headers["Cc"] = config.GmailAddress
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=\"%s\"", boundary)

	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")

	// HTML body
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	message.WriteString("\r\n")
	
	htmlBody := fmt.Sprintf(`<html>
<body style="font-family: Arial, sans-serif;">
<h2>Weekly Timesheet & Invoice</h2>
<p>Dear InfoObjects Team,</p>
<p>Please find attached my weekly invoice.</p>
<div style="background: #f8f9fa; padding: 20px; margin: 20px 0;">
<p><strong>Invoice:</strong> %s</p>
<p><strong>Hours:</strong> %.1f</p>
<p><strong>Amount:</strong> $%.2f</p>
</div>
<p>Best regards,<br>Dion Edge<br>%s</p>
</body>
</html>`, invoice.InvoiceNumber, invoice.TotalHours, invoice.LineTotal, config.GmailAddress)
	
	message.WriteString(htmlBody)
	message.WriteString("\r\n")

	// PDF attachment
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: application/pdf\r\n")
	message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"invoice_%s.pdf\"\r\n", invoice.InvoiceNumber))
	message.WriteString("Content-Transfer-Encoding: base64\r\n")
	message.WriteString("\r\n")
	
	// Base64 encode PDF
	encoded := encodeBase64(pdfData)
	message.WriteString(encoded)
	message.WriteString("\r\n")
	
	message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	// Send via SMTP
	auth := smtp.PlainAuth("", config.GmailAddress, config.GmailPassword, config.SMTPHost)
	
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         config.SMTPHost,
	}

	conn, err := tls.Dial("tcp", config.SMTPHost+":465", tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS dial failed: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, config.SMTPHost)
	if err != nil {
		return fmt.Errorf("SMTP client failed: %v", err)
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %v", err)
	}

	if err = client.Mail(config.GmailAddress); err != nil {
		return fmt.Errorf("SMTP MAIL failed: %v", err)
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("SMTP RCPT failed: %v", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA failed: %v", err)
	}

	_, err = w.Write([]byte(message.String()))
	if err != nil {
		return fmt.Errorf("SMTP write failed: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("SMTP close failed: %v", err)
	}

	return nil
}

// encodeBase64 encodes data to base64 with line breaks
func encodeBase64(data []byte) string {
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result strings.Builder
	
	for i := 0; i < len(data); i += 3 {
		b := (data[i] & 0xFC) >> 2
		result.WriteByte(base64Table[b])
		
		b = (data[i] & 0x03) << 4
		if i+1 < len(data) {
			b |= (data[i+1] & 0xF0) >> 4
			result.WriteByte(base64Table[b])
			
			b = (data[i+1] & 0x0F) << 2
			if i+2 < len(data) {
				b |= (data[i+2] & 0xC0) >> 6
				result.WriteByte(base64Table[b])
				b = data[i+2] & 0x3F
				result.WriteByte(base64Table[b])
			} else {
				result.WriteByte(base64Table[b])
				result.WriteByte('=')
			}
		} else {
			result.WriteByte(base64Table[b])
			result.WriteString("==")
		}
		
		if (i+3)%57 == 0 {
			result.WriteString("\r\n")
		}
	}
	
	return result.String()
}
