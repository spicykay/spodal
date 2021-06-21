package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var auth spotify.Authenticator

func authenticateApplication() {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}
	client := spotify.Authenticator{}.NewClient(token)
	msg, page, err := client.FeaturedPlaylists()
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	fmt.Println(msg)
	for _, playlist := range page.Playlists {
		fmt.Println("  ", playlist.Name)
	}
}

func authenticateUser(c *gin.Context) {
	redirectURL := "http://localhost:8080/authentication_callback"
	authSessionID := uuid.NewString()
	auth = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate)
	clientId := os.Getenv("SPOTIFY_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")

	auth.SetAuthInfo(clientId, clientSecret)

	// get the user to this URL - how you do that is up to you
	// you should specify a unique state string to identify the session
	url := auth.AuthURL(authSessionID)
	fmt.Printf("url %s", url)
	c.JSON(200, gin.H{
		"redirectURI": url,
	})
}

/*
https://example.com/callback?code=NApCCg..BkWtQ&state=profile%2Factivity
*/
func authenticationCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	// Error Validation
	respondWithError := func() {
		c.JSON(400, gin.H{
			"error": "Something went wrong",
		})
	}

	if code == "" || state == "" {
		respondWithError()
	}

	token, err := auth.Token(state, c.Request)
	if err != nil {
		respondWithError()
	}

	client := auth.NewClient(token)

	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	c.JSON(200, gin.H{
		"code":      code,
		"state":     state,
		"playlists": playlists.Playlists,
	})
}

func main() {
	r := gin.Default()
	r.GET("/authenticate_user", authenticateUser)
	r.GET("/authentication_callback", authenticationCallbackHandler)
	r.Run()
}
