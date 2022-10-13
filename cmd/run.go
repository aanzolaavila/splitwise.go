package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aanzolaavila/splitwise.go"
)

func main() {
	token := os.Getenv("TOKEN")

	client := splitwise.Client{
		Token: token,
	}

	ctx := context.Background()

	user, err := client.GetCurrentUser(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	userId := user.ID

	user, err = client.GetUser(ctx, userId)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("user: %+v\n", user)

	params := splitwise.UserParams{}
	originalName := user.FirstName
	params[splitwise.UserFirstname] = "Alexander"
	user, err = client.UpdateUser(ctx, userId, params)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("user: %+v\n", user)

	params = splitwise.UserParams{}
	params["first_name"] = originalName
	user, err = client.UpdateUser(ctx, userId, params)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("user: %+v\n", user)

	// ---

	groups, err := client.GetGroups(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("# groups: %d\n", len(groups))
	fmt.Printf("groups: %+v\n", groups)
}
