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
		// get code from query string
		code := r.URL.Query().Get("code")
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AccessCode, code)))
	})
}

func AddAccessTokenToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AddAccessTokenToCtx: %s", r.URL.Path)
		// get accessCode from context
		code := r.Context().Value(AccessCode).(string)
		if code == "" {
			// log.Println("No access code found in context")
			h.ServeHTTP(w, r)
		}
		// get access token from cognito
		token, err := auth.GetAccessJWT(code)
		if err != nil {
			if strings.Contains(err.Error(), "invalid_grant") {
				// redirect to login page
				http.Redirect(w, r, conf.Auth.SignInUrl, http.StatusFound)
			}
		}
		// set access token in context
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AccessToken, token)))
	})
}

func AddClientJWTStringToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// extract JWT from request
		// jwtToken := r.Header.Get("Authorization")
		// extract JWT from cookie
		jwtToken, err := r.Cookie("jwt")
		if err != nil {
			log.Println("No JWT cookie found")
			h.ServeHTTP(w, r)
			return
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), JWTEncodedString, jwtToken.Value)))
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

func AuthenticateJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get JWT in context
		if r.Context().Value(JWTEncodedString) == nil {
			// log.Println("no JWT to authenticate")
			// error
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
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if jwtToken.Expiration().Before(time.Now()) || jwt.Validate(jwtToken) != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			// cancel context
			r.Context().Done()
			return
		}
		// jwt claims
		claims, _ := jwtToken.AsMap(r.Context())

		log.Printf("claims: %v", claims)
		// Token is authenticated, pass it through
		h.ServeHTTP(w, r)
	})
}
