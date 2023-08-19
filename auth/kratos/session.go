package kratos

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"

	"log/slog"

	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/joshuaschlichting/gocms/config"
	"golang.org/x/net/publicsuffix"
)

const KRATOS_WHOAMI = "http://kratos:4433/sessions/whoami"

var conf *config.Config

var logger *slog.Logger

func InitKratos(c *config.Config) {
	conf = c
}
func SetLogger(l *slog.Logger) {
	logger = l
}

// GetJWT performs the whoami request to kratos and then constructs a JWT from the response
func GetJWT(cookie string) (string, error) {
	// Create an HTTP client with a cookie jar
	// The cookie jar is used to handle cookies automatically
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Println("Error creating cookie jar:", err)
		return "", err
	}
	client := &http.Client{
		Jar: jar,
	}
	// Create the GET request
	req, err := http.NewRequest("GET", KRATOS_WHOAMI, nil)
	req.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session", // replace with the actual cookie name
		Value: cookie,
		Path:  "/", // usually, the path is root
	})

	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	// log the entirety of the GET request
	logger.Debug("kratos request", "request", req)

	// // Set the necessary headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "http://web:8000")

	// // Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching %s: %v ", KRATOS_WHOAMI, err)
		return "", err
	}
	defer resp.Body.Close()

	// // Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}
	logger.Debug("kratos response", "body", string(body))
	// Unmarshal the response JSON (if you have a struct to unmarshal into)
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return "", err
	}
	if data["error"] != nil {
		logger.Error("kratos response error", "error", data["error"])
		return "", err
	}

	data["authSource"] = "kratos"
	var tokenString string
	//	build a JWT string from DATA
	tokenAuth := jwtauth.New("HS256", []byte(conf.Auth.JWT.SecretKey), nil)
	// conf.Auth.JWT.ExpirationTime add to now
	expirationTime := time.Now().Add(time.Second * time.Duration(conf.Auth.JWT.ExpirationTime))

	// "username": data["identity"].(map[string]interface{})["traits"].(map[string]interface{})["email"],
	log.Println("data:", data)
	// set token expiration
	_, tokenString, _ = tokenAuth.Encode(map[string]interface{}{
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
		"iss": conf.Auth.JWT.Issuer + "-kratos",
		"aud": conf.Auth.JWT.Audience,
		"sub": conf.Auth.JWT.Subject,
		// guid for jti
		"jti": uuid.New().String(),
	})

	return tokenString, nil
}

func GetUserInfo(jwt string) (string, string, error) {

	return "", "", nil
}
