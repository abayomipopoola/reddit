package web

import (
	"net/http"

	"github.com/abayomipopoola/reddit"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	store    reddit.Store
	sessions *scs.SessionManager
}

func (h *UserHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = UserRegister(w, UserParams{
			SessionData: GetSessionData(h.sessions, r.Context()),
			CSRF:        csrf.TemplateField(r),
		})
	}
}

func (h *UserHandler) RegisterSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := RegisterForm{
			Username:      r.FormValue("username"),
			Password:      r.FormValue("password"),
			UsernameTaken: false,
		}
		if _, err := h.store.UserByUsername(form.Username); err == nil {
			form.UsernameTaken = true
		}
		if !form.Validate() {
			h.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := h.store.CreateUser(&reddit.User{
			ID:       uuid.New(),
			Username: form.Username,
			Password: string(password),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.sessions.Put(r.Context(), "flash", "Your registration was successful. Please log in.")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h *UserHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = UserLogin(w, UserParams{
			SessionData: GetSessionData(h.sessions, r.Context()),
			CSRF:        csrf.TemplateField(r),
		})
	}
}

func (h *UserHandler) LoginSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := LoginForm{
			Username:             r.FormValue("username"),
			Password:             r.FormValue("password"),
			IncorrectCredentials: false,
		}
		user, err := h.store.UserByUsername(form.Username)
		if err != nil {
			form.IncorrectCredentials = true
		} else {
			compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
			form.IncorrectCredentials = compareErr != nil
		}
		if !form.Validate() {
			h.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}

		h.sessions.Put(r.Context(), "user_id", user.ID)
		h.sessions.Put(r.Context(), "flash", "You have been logged in sucessfully.")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h *UserHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.sessions.Remove(r.Context(), "user_id")
		h.sessions.Put(r.Context(), "flash", "You have been logged out sucessfully.")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
