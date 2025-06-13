package helper

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/jeffrpowell/listaway/internal/constants"
)

func SendEmailOverSMTP(email string, emailSubject string, plainBody string) (error, bool) {
	// If SMTP is not configured, just log the email
	if constants.SMTP_HOST == "" || constants.SMTP_USER == "" {
		log.Printf("[EMAIL STUB] To: %s, Subject: %s\nBody:\n%s", email, emailSubject, plainBody)
		return nil, true
	}

	mimeHeaders := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"From: " + constants.SMTP_FROM + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: " + emailSubject + "\r\n\r\n" +
		plainBody

	// Determine port
	smtpPort, _ := strconv.Atoi(constants.SMTP_PORT)
	if smtpPort == 0 {
		smtpPort = 587 // Default to 587 if port is invalid
	}

	// Get server address
	smtpAddr := fmt.Sprintf("%s:%d", constants.SMTP_HOST, smtpPort)
	to := []string{email}

	// Check if secure connection is required (default is true)
	secureConnection := true
	if strings.ToLower(constants.SMTP_SECURE) == "false" {
		secureConnection = false
	}

	var err error
	// For port 587, we need to use STARTTLS instead of direct TLS
	if secureConnection {
		if smtpPort == 587 { // SMTP with STARTTLS
			// Connect to the server
			c, err := smtp.Dial(smtpAddr)
			if err != nil {
				log.Printf("Failed to connect to SMTP server: %v", err)
				return nil, true
			}
			defer c.Close()

			// Set TLS config
			tlsConfig := &tls.Config{
				ServerName: constants.SMTP_HOST,
			}

			// Start TLS
			if err = c.StartTLS(tlsConfig); err != nil {
				log.Printf("Failed to start TLS: %v", err)
				return nil, true
			}

			// Auth
			auth := smtp.PlainAuth("", constants.SMTP_USER, constants.SMTP_PASSWORD, constants.SMTP_HOST)
			if err = c.Auth(auth); err != nil {
				log.Printf("SMTP authentication failed: %v", err)
				return nil, true
			}

			// Set the sender and recipient
			if err = c.Mail(constants.SMTP_FROM); err != nil {
				log.Printf("Failed to set sender: %v", err)
				return nil, true
			}

			if err = c.Rcpt(email); err != nil {
				log.Printf("Failed to set recipient: %v", err)
				return nil, true
			}

			// Send the email body
			w, err := c.Data()
			if err != nil {
				log.Printf("Failed to open data writer: %v", err)
				return nil, true
			}
			_, err = w.Write([]byte(mimeHeaders))
			if err != nil {
				log.Printf("Failed to write email content: %v", err)
				return nil, true
			}
			err = w.Close()
			if err != nil {
				log.Printf("Failed to close data writer: %v", err)
				return nil, true
			}

			// Send the QUIT command and close the connection
			err = c.Quit()
			if err != nil {
				log.Printf("Failed to quit SMTP session: %v", err)
				return nil, true
			}
		} else if smtpPort == 465 { // SMTP over SSL (direct TLS)
			// Create TLS config
			tlsConfig := &tls.Config{
				ServerName: constants.SMTP_HOST,
			}

			// Connect to the SMTP Server over TLS
			conn, err := tls.Dial("tcp", smtpAddr, tlsConfig)
			if err != nil {
				log.Printf("Failed to establish TLS connection: %v", err)
				return nil, true
			}
			defer conn.Close()

			// Create a client
			c, err := smtp.NewClient(conn, constants.SMTP_HOST)
			if err != nil {
				log.Printf("Failed to create SMTP client: %v", err)
				return nil, true
			}
			defer c.Close()

			// Auth
			auth := smtp.PlainAuth("", constants.SMTP_USER, constants.SMTP_PASSWORD, constants.SMTP_HOST)
			if err = c.Auth(auth); err != nil {
				log.Printf("SMTP authentication failed: %v", err)
				return nil, true
			}

			// Set the sender and recipient
			if err = c.Mail(constants.SMTP_FROM); err != nil {
				log.Printf("Failed to set sender: %v", err)
				return nil, true
			}

			if err = c.Rcpt(email); err != nil {
				log.Printf("Failed to set recipient: %v", err)
				return nil, true
			}

			// Send the email body
			w, err := c.Data()
			if err != nil {
				log.Printf("Failed to open data writer: %v", err)
				return nil, true
			}
			_, err = w.Write([]byte(mimeHeaders))
			if err != nil {
				log.Printf("Failed to write email content: %v", err)
				return nil, true
			}
			err = w.Close()
			if err != nil {
				log.Printf("Failed to close data writer: %v", err)
				return nil, true
			}

			// Send the QUIT command and close the connection
			err = c.Quit()
			if err != nil {
				log.Printf("Failed to quit SMTP session: %v", err)
				return nil, true
			}
		} else { // For other ports, try the high-level SendMail with TLS
			auth := smtp.PlainAuth("", constants.SMTP_USER, constants.SMTP_PASSWORD, constants.SMTP_HOST)
			err = smtp.SendMail(smtpAddr, auth, constants.SMTP_FROM, to, []byte(mimeHeaders))
			if err != nil {
				log.Printf("Failed to send email with TLS: %v", err)
				return nil, true
			}
		}
	} else {
		// Use standard SMTP without TLS
		auth := smtp.PlainAuth("", constants.SMTP_USER, constants.SMTP_PASSWORD, constants.SMTP_HOST)
		err = smtp.SendMail(smtpAddr, auth, constants.SMTP_FROM, to, []byte(mimeHeaders))
	}
	return err, false
}
