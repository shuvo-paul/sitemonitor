package handler

import (
	"fmt"
	"net/http"

	"github.com/shuvo-paul/uptimebot/internal/auth/model"
	"github.com/shuvo-paul/uptimebot/internal/auth/service"
	"github.com/shuvo-paul/uptimebot/internal/renderer"
	"github.com/shuvo-paul/uptimebot/pkg/flash"
)

type UserHandler struct {
	Template struct {
		Register *renderer.Template
		Login    *renderer.Template
	}
	sessionService service.SessionServiceInterface
	authService    service.AuthServiceInterface
	flashStore     flash.FlashStoreInterface
}

func NewUserHandler(
	authService service.AuthServiceInterface,
	sessionService service.SessionServiceInterface,
	flashStore flash.FlashStoreInterface,
) *UserHandler {
	return &UserHandler{
		authService:    authService,
		sessionService: sessionService,
		flashStore:     flashStore,
	}
}

func (c *UserHandler) ShowRegisterForm(w http.ResponseWriter, r *http.Request) {
	if c.redirectIfAuthenticated(w, r) {
		return
	}

	flashId := flash.GetFlashIDFromContext(r.Context())

	data := map[string]any{
		"Title": "Registration",
		"Erros": c.flashStore.GetFlash(flashId, "errors"),
	}

	c.Template.Register.Render(w, r, data)
}

func (c *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	user := &model.User{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	_, err := c.authService.CreateUser(user)
	if err != nil {
		flashId := flash.GetFlashIDFromContext(r.Context())

		errors := []string{err.Error()}
		fmt.Println(errors[0])
		c.flashStore.SetFlash(flashId, "errors", errors)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (c *UserHandler) ShowLoginForm(w http.ResponseWriter, r *http.Request) {
	if c.redirectIfAuthenticated(w, r) {
		return
	}

	data := map[string]string{
		"Title": "Login",
	}
	c.Template.Login.Render(w, r, data)
}

func (c *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := c.authService.Authenticate(email, password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, token, err := c.sessionService.CreateSession(user.ID)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/targets", http.StatusSeeOther)
}

func (c *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Invalidate the session in the backend
		if err := c.sessionService.DeleteSession(cookie.Value); err != nil {
			http.Error(w, "Failed to logout", http.StatusInternalServerError)
			return
		}
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (c *UserHandler) redirectIfAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	if cookie, err := r.Cookie("session_token"); err == nil {
		user, err := c.sessionService.ValidateSession(cookie.Value)

		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
			return false
		}
		if user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return true
		}
	}
	return false
}
