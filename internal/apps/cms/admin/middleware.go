package admin

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	auth "github.com/joshuaschlichting/gocms/auth/oauth2"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/internal/apps/cms/data/db"
	"github.com/lestrrat-go/jwx/jwt"
)

type MiddlewareWithDB interface {
	AddUserToCtx(h http.Handler) http.Handler
}
type contextKey int

const (
	AccessCode contextKey = iota
	AccessToken
	User
	JWTEncodedString
)

func NewMiddlewareWithDB(db db.DBCache, jwtSecretKey string) DBMiddleware {
	r := DBMiddleware{
		db:           db,
		jwtSecretKey: jwtSecretKey,
	}
	return r
}

type DBMiddleware struct {
	db           db.DBCache
	jwtSecretKey string
}

func (m DBMiddleware) AddUserToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt := r.Context().Value(JWTEncodedString).(string)
		// unmarshal jwt

		// cast to *jwtauth.JWTAuth
		tokenAuth := jwtauth.New("HS256", []byte(m.jwtSecretKey), nil)
		jwtAuthToken, err := tokenAuth.Decode(jwt)
		if err != nil {
			logger.Error(fmt.Sprintf("error decoding JWT from context: %s: JWT from context: %v", err.Error(), jwt))
			h.ServeHTTP(w, r)
			return
		}
		username, _ := jwtAuthToken.Get("userInfo")
		user, err := m.db.GetUserByName(r.Context(), username.(string))
		if err != nil {
			logger.Error(fmt.Sprintf("unable to add user to context in middleware: unable to get user '%s' from db due to error: %v", username, err.Error()))
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), User, user)))
	})
}

// AuthorizeUserGroup checks if the user is in the specified group
func (m DBMiddleware) AuthorizeUserGroup(group string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user from context
			user := r.Context().Value(User).(db.User)
			if user.Name == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			userInGroup, err := m.db.GetUserIsInGroup(r.Context(), db.GetUserIsInGroupParams{
				UserID:        user.ID,
				UsergroupName: group,
			})
			if err != nil {
				logger.Error("unable to check if user '%s' is in group '%s' due to error: %v", user.Name, group, err.Error())
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if !userInGroup {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			// pass through
			h.ServeHTTP(w, r)
		})
	}
}

var conf *config.Config

func InitMiddleware(config *config.Config) {
	conf = config
}

func init() {
	// find and load config.yml independently of the app

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
			logger.Debug("No JWT token found in header")
		}
		cookieJWTToken, err := r.Cookie("jwt")
		if err != nil {
			log.Println("No JWT cookie found")
			h.ServeHTTP(w, r)
			return
		}
		logger.Debug("JWT found in cookie")
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

		logger.Debug("JWT claims", "JWT", claims)
		h.ServeHTTP(w, r)
	})
}
