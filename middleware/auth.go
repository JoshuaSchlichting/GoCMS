package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/joshuaschlichting/gocms/auth"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/db"
	"github.com/lestrrat-go/jwx/jwt"
)

var conf *config.Config

func InitMiddleware(config *config.Config) {
	conf = config
}

func AddURLAccessCodeToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AccessCode, code)))
	})
}

func AddAccessTokenToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AddAccessTokenToCtx: %s", r.URL.Path)
		code := r.Context().Value(AccessCode).(string)
		if code == "" {
			h.ServeHTTP(w, r)
		}
		token, err := auth.GetAccessJWT(code)
		if err != nil {
			if strings.Contains(err.Error(), "invalid_grant") {
				http.Redirect(w, r, conf.Auth.SignInUrl, http.StatusFound)
			}
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AccessToken, token)))
	})
}

func AddClientJWTStringToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		jwtToken := r.Header.Get("Authorization")
		if jwtToken == "" {
			log.Println("No JWT token found in header")
		}
		cookieJWTToken, err := r.Cookie("jwt")
		if err != nil {
			log.Println("No JWT cookie found")
			h.ServeHTTP(w, r)
			return
		}
		log.Println("JWT found in cookie")
		jwtToken = cookieJWTToken.Value

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), JWTEncodedString, jwtToken)))
	})
}

func AddNewJwtToCtxCookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(User) == nil {
			h.ServeHTTP(w, r)
			return
		}
		user := r.Context().Value(User).(db.User)
		if user.Name == "" {
			log.Println("user was set in context but is empty string")
			h.ServeHTTP(w, r)
		}
		tokenAuth := jwtauth.New("HS256", []byte(conf.Auth.JWT.SecretKey), nil)
		// conf.Auth.JWT.ExpirationTime add to now
		expirationTime := time.Now().Add(time.Second * time.Duration(conf.Auth.JWT.ExpirationTime))
		// set token expiration

		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
			"username": user.Name,
			"exp":      expirationTime.Unix(),
			"iat":      time.Now().Unix(),
			"iss":      conf.Auth.JWT.Issuer,
			"aud":      conf.Auth.JWT.Audience,
			"sub":      conf.Auth.JWT.Subject,
			// guid for jti
			"jti": uuid.New().String(),
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: tokenString,
		})
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), JWTEncodedString, tokenString)))
	})
}

// AuthenticateJWT checks if the JWT is valid via the secret key and expiration date
func AuthenticateJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get JWT in context
		if r.Context().Value(JWTEncodedString) == nil {
			log.Println("unauthorized: no JWT to authenticate")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		token := r.Context().Value(JWTEncodedString).(string)
		if token == "" {
			// log.Println("No JWT token found in context")
			h.ServeHTTP(w, r)
			return
		}
		// cast to *jwtauth.JWTAuth
		tokenAuth := jwtauth.New("HS256", []byte(conf.Auth.JWT.SecretKey), nil)
		jwtToken, err := tokenAuth.Decode(token)
		if err != nil {
			log.Println("unauthorized: invalid JWT token (perhaps the secret key has changed?)")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			r.Context().Done()
			return
		}
		if jwtToken.Expiration().Before(time.Now()) {
			log.Println("unauthorized: JWT token has expired")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			r.Context().Done()
			return
		}
		if jwt.Validate(jwtToken) != nil {
			log.Println("unauthorized: JWT claims are invalid")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			r.Context().Done()
			return
		}

		claims, _ := jwtToken.AsMap(r.Context())

		log.Printf("claims: %v", claims)
		h.ServeHTTP(w, r)
	})
}
