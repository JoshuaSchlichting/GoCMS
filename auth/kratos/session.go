package kratos

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
)

// GetJWT performs the whoami request to kratos and then constructs a JWT from the response
func GetJWT(cookie string) (string, error) {
	// Create an HTTP client with a cookie jar
	// The cookie jar is used to handle cookies automatically
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	// Set the request URL
	url := "http://kratos:4433/sessions/whoami"

	// Create the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Set the necessary headers
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching Who am I?:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	// Unmarshal the response JSON (if you have a struct to unmarshal into)
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return "", err
	}
	data["authSource"] = "kratos"

	// Print the data
	fmt.Println(data)
	return "", nil
}

func GetUserInfo(jwt string) (string, string, error) {
	return "", "", nil
}
