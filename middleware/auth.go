package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/joshuaschlichting/gocms/auth"
	"github.com/joshuaschlichting/gocms/config"
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
		// get accessCode from context
		code := r.Context().Value(AccessCode).(string)
		if code == "" {
			// log.Println("No access code found in context")
			h.ServeHTTP(w, r)
		}
		// get access token from cognito
		token, err := getAccessToken(code)
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

func AddUserInfoToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get access token in context
		if r.Context().Value(AccessToken) == nil {
			// log.Println("No access token found in context")
			h.ServeHTTP(w, r)
			return
		}
		// get access token from context
		token := r.Context().Value(AccessToken).(string)
		if token == "" {
			// log.Println("No access token found in context")
			h.ServeHTTP(w, r)
		}
		// get user info from cognito
		cognitoClient, _ := auth.NewCognito()

		userInfo, _, _ := cognitoClient.GetUserInfo(token)
		// set user info in context
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserInfo, userInfo)))
	})
}

func AddClientJWTToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// extract JWT from request
		jwtToken := r.Header.Get("Authorization")
		if jwtToken == "" {
			// log.Println("No JWT token found in context")
			h.ServeHTTP(w, r)
			return
		}
		// cast to *jwtauth.JWTAuth
		// tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		// jwtAuthToken, _ := tokenAuth.Decode(jwtToken)
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), JWTToken, jwtToken)))
	})
}

func AddNewJwtToCtxCookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(UserInfo) == nil {
			h.ServeHTTP(w, r)
			return
		}
		userInfo := r.Context().Value(UserInfo).(string)
		if userInfo == "" {
			log.Println("user was set in context but is empty string")
			h.ServeHTTP(w, r)
		}
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		expirationTime := time.Now().Add(5 * time.Minute)

		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
			"userInfo": userInfo,
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
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), JWTToken, tokenString)))
	})
}

func AuthenticateJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get JWT in context
		if r.Context().Value(JWTToken) == nil {
			// log.Println("no JWT to authenticate")
			h.ServeHTTP(w, r)
			return
		}
		token := r.Context().Value(JWTToken).(string)
		if token == "" {
			// log.Println("No JWT token found in context")
			h.ServeHTTP(w, r)
			return
		}
		// cast to *jwtauth.JWTAuth
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		jwtToken, err := tokenAuth.Decode(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// get jwtToken is expired

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

func getAccessToken(authorizationCode string) (string, error) {
	if authorizationCode == "" {
		return "", errors.New("no authorization code found")
	}
	cognitoClient, _ := auth.NewCognito()
	payload, err := cognitoClient.GetCognitoTokenEndpointPayload(authorizationCode)
	if err != nil {
		log.Printf("Error getting token endpoint payload: %v\n", err)
		return "", err
	}
	return payload.AccessToken, nil
}
