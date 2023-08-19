package kratos

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

func GetUserInfo(orySessionCookie string) (string, string, error) {
	logger.Debug("kratos.GetUserInfo ->")

	data, err := getOryPayload(orySessionCookie)
	if err != nil {
		logger.Error("error getting ory payload", "error", err)
		return "", "", err
	}

	username := data["identity"].(map[string]interface{})["traits"].(map[string]interface{})["username"].(string)
	email := data["identity"].(map[string]interface{})["traits"].(map[string]interface{})["email"].(string)

	logger.Debug("kratos.GetUserInfo <-", "username", username, "email", email)
	return username, email, nil
}

func getOryPayload(cookie string) (map[string]interface{}, error) {
	// Create an HTTP client with a cookie jar
	// The cookie jar is used to handle cookies automatically
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Println("Error creating cookie jar:", err)
		return map[string]interface{}{}, err
	}
	client := &http.Client{
		Jar: jar,
	}
	// Create the GET request
	req, err := http.NewRequest("GET", kratos_whoami_url, nil)
	req.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session", // replace with the actual cookie name
		Value: cookie,
		Path:  "/", // usually, the path is root
	})

	if err != nil {
		fmt.Println("Error creating request:", err)
		return map[string]interface{}{}, err
	}
	// log the entirety of the GET request
	logger.Debug("GET to kratos", "request", req)

	// // Set the necessary headers
	req.Header.Set("Accept", "application/json")

	// // Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching %s: %v ", kratos_whoami_url, err)
		logger.Error(fmt.Sprint("error fetching ", kratos_whoami_url), "error", err)
		return map[string]interface{}{}, err
	}
	defer resp.Body.Close()

	// // Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return map[string]interface{}{}, err
	}
	logger.Debug("kratos response", "body", string(body))
	// Unmarshal the response JSON (if you have a struct to unmarshal into)
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return map[string]interface{}{}, err
	}
	logger.Debug("kratos response", "data", data)
	return data, nil
}
