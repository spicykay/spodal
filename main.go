package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	fmt.Println("Hello, World!")
	clientId := os.Getenv("SPOTIFY_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")
	fmt.Printf("ID: %s, Secret: %s", clientId, clientSecret)

	config := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
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
