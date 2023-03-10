package cookieservice

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

type cookieServiceImpl struct {
	secretKey []byte
}

func New(key string) (CookieService, error) {
	return &cookieServiceImpl{[]byte(key)}, nil
}

func (c *cookieServiceImpl) AuthenticateUser(w http.ResponseWriter, r *http.Request) (string, error) {
	userID, authErr := c.ReadSigned(r, "userID")
	if authErr != nil {
		return "", authErr
	}
	return userID, nil
}

func (c *cookieServiceImpl) WriteSigned(w http.ResponseWriter, userID string) error {
	cookie := http.Cookie{
		Name:     "userID",
		Value:    userID,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
	}
	mac := hmac.New(sha256.New, c.secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := mac.Sum(nil)

	cookie.Value = string(signature) + cookie.Value

	return write(w, cookie)
}

func (c *cookieServiceImpl) ReadSigned(r *http.Request, name string) (string, error) {
	// {signature}{original value}
	signedValue, err := read(r, name)
	if err != nil {
		return "", err
	}

	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}

	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]
	mac := hmac.New(sha256.New, c.secretKey)
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	return value, nil
}

func write(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(cookie.Value))
	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}
	http.SetCookie(w, &cookie)
	return nil
}

func read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	value, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}
	return string(value), nil
}
