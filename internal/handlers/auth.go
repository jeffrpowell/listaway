package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	constants.ROUTER.HandleFunc("/", middleware.DefaultMiddlewareChain(rootHandler)).Methods("GET")
	constants.ROUTER.HandleFunc("/auth", middleware.DefaultPublicMiddlewareChain(authHandler))
	constants.ROUTER.HandleFunc("/reset", middleware.DefaultPublicMiddlewareChain(resetHandler)).Methods("POST")
	constants.ROUTER.HandleFunc("/reset/{token}", middleware.DefaultPublicMiddlewareChain(resetTokenHandler))
	constants.ADMIN_EXISTS = database.AdminUserExists()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		authGET(w, r)
	case "POST":
		authPOST(w, r)
	case "DELETE":
		authDELETE(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

/* Password reset POST */
func resetHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	// Check if user exists with this email
	userID, err := database.GetUserByEmail(email)
	if err != nil {
		log.Printf("Error checking user email: %v", err)
		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
		return
	}

	// Always return success even if email not found (security best practice)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("If your email is registered, a reset link has been sent."))

	// Only proceed if user was found
	if userID == -1 {
		return
	}

	// Generate token and store in database
	token, err := database.CreatePasswordResetToken(email)
	if err != nil {
		log.Printf("Error creating reset token: %v", err)
		return
	}

	// Send reset email
	sendResetEmail(email, token)
}

// Helper: send reset email with SMTP server if configured
func sendResetEmail(email, token string) {
	// Get application base URL from configuration
	baseURL := constants.APP_URL
	resetURL := fmt.Sprintf("%s/reset/%s", baseURL, token)

	// Email template
	emailSubject := "Password Reset Request"
	// Build plain text email body
	plainBody := fmt.Sprintf(
		"Hello,\n\nA password reset was requested for your Listaway account.\n\n"+
			"Please click the following link to reset your password:\n%s\n\n"+
			"This link will expire in 1 hour.\n\n"+
			"If you did not request this password reset, please ignore this email.\n\n"+
			"Regards,\nThe Listaway Team", resetURL)

	// If SMTP is not configured, just log the email
	if constants.SMTP_HOST == "" || constants.SMTP_USER == "" {
		log.Printf("[EMAIL STUB] To: %s, Subject: %s\nBody:\n%s", email, emailSubject, plainBody)
		return
	}

	// Build MIME email with headers
	mimeHeaders := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"From: " + constants.SMTP_FROM + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: " + emailSubject + "\r\n\r\n" +
		plainBody

	// Set up authentication
	auth := smtp.PlainAuth("", constants.SMTP_USER, constants.SMTP_PASSWORD, constants.SMTP_HOST)

	// Determine port
	smtpPort, _ := strconv.Atoi(constants.SMTP_PORT)
	if smtpPort == 0 {
		smtpPort = 587 // Default to 587 if port is invalid
	}

	// Send email
	smtpAddr := fmt.Sprintf("%s:%d", constants.SMTP_HOST, smtpPort)
	to := []string{email}

	err := smtp.SendMail(smtpAddr, auth, constants.SMTP_FROM, to, []byte(mimeHeaders))
	if err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		return
	}

	log.Printf("Password reset email sent to %s", email)
}

/* Password Reset page */
func resetTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	token = strings.TrimSpace(token)

	// Validate token from database
	email, valid, err := database.ValidatePasswordResetToken(token)
	if err != nil {
		log.Printf("Error validating reset token: %v", err)
		web.ResetFormPage(w, false)
		return
	}

	if !valid {
		web.ResetFormPage(w, false)
		return
	}

	switch r.Method {
	case "GET":
		web.ResetFormPage(w, true)
	case "POST":
		password := r.FormValue("password")
		err := database.UpdateUserPassword(email, password)
		if err != nil {
			log.Printf("Error updating user password: %v", err)
			http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
			return
		}

		_ = database.InvalidatePasswordResetToken(token)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Password updated. You may now log in."))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

/* Login page */
func authGET(w http.ResponseWriter, r *http.Request) {
	web.LoginPage(w)
}

/* Login */
func authPOST(w http.ResponseWriter, r *http.Request) {
	session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)

	userId, err := database.LoginUser(r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if userId == -1 {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Values["userId"] = userId
	session.Save(r, w)
	w.Header().Add("Location", "/list")
	w.WriteHeader(http.StatusOK)
}

/* Logout */
func authDELETE(w http.ResponseWriter, r *http.Request) {
	session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)

	// Revoke users authentication
	session.Values["authenticated"] = false
	delete(session.Values, "userId")
	session.Options.MaxAge = -1
	session.Save(r, w)
	w.Header().Add("Location", "/auth")
	w.WriteHeader(http.StatusNoContent)
}
