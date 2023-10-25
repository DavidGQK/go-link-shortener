package server

import (
	"context"
	"fmt"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/storage/initstorage"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const SecretKey = "VerySecretKey"

func (s *Server) CookieMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		cookie, err := r.Cookie("shortener_session")
		if err != nil {
			newCookie, err := createNewCookie(s.storage)
			if err != nil {
				logger.Log.Error("create cookie error", zap.Error(err))
				http.Error(w, "Internal Backend Error", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "shortener_session",
				Value: newCookie,
			})
			r.AddCookie(&http.Cookie{
				Name:  "shortener_session",
				Value: newCookie,
			})
		} else if !isCookieValid(cookie.Value, s.storage) {
			if path == "/api/user/urls" {
				logger.Log.Error("invalid cookie", zap.Error(err))
				http.Error(w, "invalid cookie", http.StatusUnauthorized)
				return
			}

			newCookie, err := createNewCookie(s.storage)
			if err != nil {
				logger.Log.Error("create cookie error", zap.Error(err))
				http.Error(w, "Internal Backend Error", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "shortener_session",
				Value: newCookie,
			})
			r.AddCookie(&http.Cookie{
				Name:  "shortener_session",
				Value: newCookie,
			})
		}
		h.ServeHTTP(w, r)
	}
}

func createNewCookie(rep repository) (cookie string, err error) {
	if rep.GetMode() == initstorage.DBMode {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		user, err := rep.CreateUser(ctx)
		if err != nil {
			return "", err
		}
		userID := user.UserID

		cookie, err = BuildJWTString(userID)
		if err != nil {
			return cookie, err
		}
		err = rep.UpdateUser(ctx, userID, cookie)
		if err != nil {
			return cookie, err
		}
	} else {
		logger.Log.Infow("not database mode")
		userID := 1
		cookie, err = BuildJWTString(userID)
		if err != nil {
			return cookie, err
		}
	}

	return cookie, nil
}

func BuildJWTString(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func isCookieValid(cookie string, rep repository) bool {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		logger.Log.Error("parse jwt error", zap.Error(err))
		return false
	}

	if !token.Valid {
		logger.Log.Error("token is not valid")
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	userID := claims.UserID
	_, err = rep.FindUserByID(ctx, userID)
	return err == nil
}
