package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aanzolaavila/splitwise.go"
	"golang.org/x/oauth2"
)

func getTokenClient(token string) splitwise.Client {
	return splitwise.Client{
		Token: token,
	}
}

func getOAuthClient() splitwise.Client {
	ctx := context.Background()

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://secure.splitwise.com/oauth/authorize",
			TokenURL: "https://secure.splitwise.com/oauth/token ",
		},
	}

	state := "random_string"

	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("error getting code: %v", err)
	}

	fmt.Printf("got code: %s\n", code)

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("error setting exchange: %v", err)
	}

	httpClient := conf.Client(ctx, tok)

	return splitwise.Client{
		HttpClient: httpClient,
	}
}

func main() {
	token := os.Getenv("TOKEN")
	var client splitwise.Client

	if token != "" {
		client = getTokenClient(token)
	} else {
		client = getOAuthClient()
	}

	ctx := context.Background()

	userExamples(ctx, client)
	groupExamples(ctx, client)
	friendsExamples(ctx, client)
}

func userExamples(ctx context.Context, client splitwise.Client) {
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
}

func groupExamples(ctx context.Context, client splitwise.Client) {
	groups, err := client.GetGroups(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("# groups: %d\n", len(groups))
	fmt.Printf("groups: \n")
	for i, group := range groups {
		fmt.Printf("Group #%d: %d - %s\n", i, group.ID, group.Name)
	}

	// this should fail
	const invalidGroupId = 10
	_, err = client.GetGroup(ctx, invalidGroupId)
	if err == nil {
		log.Fatalf("this should have failed")
	}
	fmt.Printf("expected error: %v\n", err)

	// Create a group
	const createGroupName = "delete this group"
	currentUser, err := client.GetCurrentUser(ctx)
	if err != nil {
		log.Fatalf("could not get current user: %v", err)
	}

	groupUser := splitwise.GroupUser{
		Id: currentUser.ID,
	}

	groupParams := splitwise.GroupParams{}
	groupParams[splitwise.GroupType] = "other"
	groupParams[splitwise.GroupSimplifyByDefault] = true
	group, err := client.CreateGroup(ctx, createGroupName, groupParams, []splitwise.GroupUser{groupUser})
	if err != nil {
		log.Fatalf("error creating group: %v", err)
	}

	fmt.Printf("Group Id %d, Name: %s\n", group.ID, group.Name)

	// let's delete it as well
	err = client.DeleteGroup(ctx, group.ID)
	if err != nil {
		log.Fatalf("error deleting group: %v", err)
	}

	fmt.Println("group deleted")

	// let's try to undelete it
	err = client.RestoreGroup(ctx, group.ID)
	if err != nil {
		log.Fatalf("error restoring group: %v", err)
	}

	fmt.Println("group restored")

	// let's delete it again
	err = client.DeleteGroup(ctx, group.ID)
	if err != nil {
		log.Fatalf("error deleting group: %v", err)
	}

	fmt.Println("group deleted again")
}

func friendsExamples(ctx context.Context, client splitwise.Client) {
	friends, err := client.GetFriends(ctx)
	if err != nil {
		log.Fatalf("error getting friends: %v", err)
	}

	fmt.Println("Friends")
	for _, f := range friends {
		fmt.Printf("friend #%d: %s\n", f.ID, f.Email)
	}

	if len(friends) > 0 {
		id := friends[0].ID
		friend, err := client.GetFriend(ctx, id)
		if err != nil {
			log.Fatalf("error getting friend #%d: %v", id, err)
		}

		fmt.Printf("friend #%d - Name %s - Email %s\n", friend.ID, friend.FirstName, friend.Email)
	}

	// add a friend
	const friendEmail = "false_friend@mail.com"
	params := splitwise.FriendParams{}
	params[splitwise.FriendFirstname] = "False"
	params[splitwise.FriendLastname] = "Friend"
	friend, err := client.AddFriend(ctx, friendEmail, params)
	if err != nil {
		log.Fatalf("failed to create friend: %v", err)
	}

	fmt.Printf("created friend: %+v\n", friend)

	friendId := friend.ID
	err = client.DeleteFriend(ctx, friendId)
	if err != nil {
		log.Fatalf("failed to delete friend: %d", friendId)
	}

	fmt.Printf("deleted friend: %d\n", friendId)
}
