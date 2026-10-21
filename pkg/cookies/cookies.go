package cookies

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound = errors.New("cookie not found")
)

// Options — общие атрибуты, применяемые ко всем cookie,
// которые ставит/удаляет пакет. Задаются один раз при инициализации.
type Options struct {
	Path     string        // по умолчанию "/"
	Domain   string        // "" — текущий домен
	Secure   bool          // true в проде (HTTPS)
	HttpOnly bool          // true почти всегда
	SameSite http.SameSite // по умолчанию http.SameSiteStrictMode
}

type manager struct {
	opts Options
}

type Manager interface {
	Get(r *http.Request, name string) (string, error)
	Set(w http.ResponseWriter, name, value string, maxAge int)
	Delete(w http.ResponseWriter, name string)
}

func NewManager(opts Options) *manager {
	if opts.Path == "" {
		opts.Path = "/"
	}
	if opts.SameSite == 0 {
		opts.SameSite = http.SameSiteStrictMode
	}
	return &manager{opts: opts}
}

func (m *manager) Get(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if errors.Is(err, http.ErrNoCookie) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func (m *manager) build(name, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     m.opts.Path,
		Domain:   m.opts.Domain,
		MaxAge:   maxAge,
		Secure:   m.opts.Secure,
		HttpOnly: m.opts.HttpOnly,
		SameSite: m.opts.SameSite,
	}
}

func (m *manager) Set(w http.ResponseWriter, name, value string, maxAge int) {
	if maxAge <= 0 {
		m.Delete(w, name)
		return
	}
	http.SetCookie(w, m.build(name, value, maxAge))
}

func (m *manager) Delete(w http.ResponseWriter, name string) {
	http.SetCookie(w, m.build(name, "", -1))
}
