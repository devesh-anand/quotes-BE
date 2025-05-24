package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"api.deveshanand.com/quotes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	githubClientID     string
	githubClientSecret string
	githubRedirectURL  string
)

const oauthStateString = "randomstatestring" // Should be a random, unguessable string in production

func SetupRouter() *gin.Engine {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://quotes.deveshanand.com"},
		// AllowOrigins:     []string{"http://localhost:3000", "https://quotes.deveshanand.com"},
		AllowMethods:     []string{"OPTIONS", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "X-Requested-With", "Authorization", "Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//routes
	r.GET("/", helloHandler)
	fmt.Printf("server on port 8080")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//Quotes routes
	r.GET("/quotes", quotes.QuotesHandler)
	r.GET("/authors", quotes.GetAllAuthors)

	r.POST("/quote", quotes.SubmitQuote)

	// Decap CMS OAuth routes
	r.GET("/github/auth", githubAuthHandler)
	r.GET("/github/callback", githubCallbackHandler)

	return r
}

func githubAuthHandler(c *gin.Context) {
	githubClientID = os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	githubRedirectURL = os.Getenv("REDIRECT_URL")
	// In a real application, you should generate and store a random state string
	// and verify it in the callback to prevent CSRF attacks.
	// For simplicity, we are using a constant string here.
	// Consider using a library for OAuth2 that handles this more robustly.
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=repo&state=%s",
		githubClientID,
		url.QueryEscape(githubRedirectURL), // Make sure this matches your GitHub app's callback URL
		oauthStateString,
	)
	fmt.Println(redirectURL)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func githubCallbackHandler(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		c.String(http.StatusBadRequest, "Invalid state string. CSRF Attack?")
		return
	}

	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Authorization code not found")
		return
	}

	// Exchange code for token
	tokenURL := "https://github.com/login/oauth/access_token"
	formData := url.Values{}
	formData.Set("client_id", githubClientID)
	formData.Set("client_secret", githubClientSecret)
	formData.Set("code", code)
	formData.Set("redirect_uri", githubRedirectURL)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create token request: "+err.Error())
		return
	}
	
	// GitHub expects the parameters in the request body, not the URL
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to exchange code for token: "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read token response: "+err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.String(http.StatusInternalServerError, fmt.Sprintf("GitHub API Error: %s, Body: %s", resp.Status, string(body)))
		return
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		TokenType    string `json:"token_type"`
		Error        string `json:"error"`
		ErrorDesc    string `json:"error_description"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse token response: "+err.Error())
		return
	}

	if tokenResponse.Error != "" {
		htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>OAuth Error</title>
    <script>
        if (window.opener) {
            window.opener.postMessage(
                'authorization:github:error:${JSON.stringify({ error: "%s", error_description: "%s" })}',
                '*'
            );
            setTimeout(function() {
                window.close();
            }, 1500);
        }
    </script>
</head>
<body>
    <p>Authentication Error: %s. Description: %s. Closing window...</p>
</body>
</html>`, tokenResponse.Error, tokenResponse.ErrorDesc, tokenResponse.Error, tokenResponse.ErrorDesc)

		fmt.Println(htmlContent)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
		return
	}

	if tokenResponse.AccessToken == "" {
		c.String(http.StatusInternalServerError, "Access token not found in response. Body: "+string(body))
		return
	}

	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>GitHub Authentication Success</title>
  <script>
    (function() {
      if (window.opener) {
        window.opener.postMessage(
          'authorization:github:success:${JSON.stringify({ token: "%s", provider: "github" })}',
          '*'
        );
        window.close();
      } else {
        console.warn("No window.opener detected.");
      }
    })();
  </script>
</head>
<body>
		<p>Authentication successful! You can close this window.</p>
	</body>
</html>`, tokenResponse.AccessToken)

	fmt.Println(htmlContent)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}

func helloHandler(c *gin.Context) {
	c.JSON(202, gin.H{
		"data": "Hello World",
	})

	fmt.Print("Hello World!")
}
