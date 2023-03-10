package cookieservice

import "net/http"

var _ CookieService = (*cookieServiceImpl)(nil)

type CookieService interface {
	WriteSigned(w http.ResponseWriter, userID string) error
	ReadSigned(r *http.Request, name string) (string, error)
	AuthenticateUser(w http.ResponseWriter, r *http.Request) (string, error)
}
