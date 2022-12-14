package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aanzolaavila/splitwise.go"
	"github.com/aanzolaavila/splitwise.go/resources"
	"golang.org/x/oauth2"

	gojson "github.com/goccy/go-json"
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
			TokenURL: "https://secure.splitwise.com/oauth/token",
		},
		// redirect_url is required as splitwise reports "invalid_grant" error otherwise
		RedirectURL: "https://localhost:17000/oauth/redirect",
	}

	state := "state"

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

	// let's set a third party json library for increased performance
	client.JsonMarshaler = gojson.Marshal
	client.JsonUnmarshaler = gojson.Unmarshal

	// let's use a custom logger
	logger := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)
	client.Logger = logger

	ctx := context.Background()

	userExamples(ctx, client)
	groupExamples(ctx, client)
	friendsExamples(ctx, client)
	expensesExamples(ctx, client)
	commentExamples(ctx, client)
	notificationsExamples(ctx, client)
	currenciesExamples(ctx, client)
	categoriesExamples(ctx, client)
	parsedExpenseExamples(ctx, client)
}

func userExamples(ctx context.Context, client splitwise.Client) {
	user, err := client.GetCurrentUser(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	userId := user.ID

	user, err = client.GetUser(ctx, int(userId))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("user: %+v\n", user)

	params := splitwise.UserParams{}
	originalName := user.FirstName
	params[splitwise.UserFirstname] = "Alexander"
	user, err = client.UpdateUser(ctx, int(userId), params)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("user: %+v\n", user)

	params = splitwise.UserParams{}
	params["first_name"] = originalName
	user, err = client.UpdateUser(ctx, int(userId), params)
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
	fmt.Printf("expecting NotFound error\n")
	if !errors.Is(err, splitwise.ErrNotFound) {
		log.Fatalf("error was unexpected, got [%v] and should have been [%v]", err, splitwise.ErrNotFound)
	}

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
	err = client.DeleteGroup(ctx, int(group.ID))
	if err != nil {
		log.Fatalf("error deleting group: %v", err)
	}

	fmt.Println("group deleted")

	// let's try to undelete it
	err = client.RestoreGroup(ctx, int(group.ID))
	if err != nil {
		log.Fatalf("error restoring group: %v", err)
	}

	fmt.Println("group restored")

	// let's delete it again
	err = client.DeleteGroup(ctx, int(group.ID))
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
		friend, err := client.GetFriend(ctx, int(id))
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
	err = client.DeleteFriend(ctx, int(friendId))
	if err != nil {
		log.Fatalf("failed to delete friend: %d", friendId)
	}

	fmt.Printf("deleted friend: %d\n", friendId)
}

func expensesExamples(ctx context.Context, client splitwise.Client) {
	getExpensesExample(ctx, client)
	createExpenseEqualGroupSplitExample(ctx, client)
	createExpenseBySharesExample(ctx, client)
}

func getExpensesExample(ctx context.Context, client splitwise.Client) {
	params := splitwise.ExpensesParams{}
	const monthDuration = 60 * 60 * 24 * 30
	params[splitwise.ExpensesDatedAfter] = time.Now().Add(-1 * 3 * monthDuration * time.Second)
	params[splitwise.ExpensesLimit] = 100
	expenses, err := client.GetExpenses(ctx, params)
	if err != nil {
		log.Fatalf("failed to get expenses: %v", err)
	}

	fmt.Printf("Expenses:\n")
	for _, e := range expenses {
		fmt.Printf("Expense #%d [%s]: %s\n", e.ID, e.Date, e.Description)
	}

	// Query one expense
	if len(expenses) > 0 {
		expenseId := expenses[0].ID
		expense, err := client.GetExpense(ctx, int(expenseId))
		if err != nil {
			log.Fatalf("failed to get expense #%d: %v", expenseId, err)
		}

		fmt.Printf("Expense #%d: %+v\n", expenseId, expense)
	}
}

func createExpenseEqualGroupSplitExample(ctx context.Context, client splitwise.Client) {
	// Create expense with Equal group split
	groups, err := client.GetGroups(ctx)
	if err != nil {
		log.Fatalf("could not get groups: %v", err)
	}

	if len(groups) == 1 {
		fmt.Printf("only one group (no group = 0), not doing anything\n")
		return
	}

	groupId := groups[0].ID
	groupName := groups[0].Name
	for i := 1; i < len(groups) && groupId == 0; i++ {
		groupId = groups[i].ID
		groupName = groups[i].Name
	}

	fmt.Printf("selected group: #%d - %s\n", groupId, groupName)

	// ---

	newExpenses, err := client.CreateExpenseEqualGroupSplit(ctx, 10000, "should delete", int(groupId), nil)
	if err != nil {
		log.Fatalf("could not create expense: %v", err)
	}

	fmt.Printf("expenses created: %d\n", len(newExpenses))
	for _, e := range newExpenses {
		fmt.Printf("Expense #%d: %s\n", e.ID, e.Description)
	}

	// let's delete those test expenses
	fmt.Printf("deleting created expenses\n")
	for _, e := range newExpenses {
		err := client.DeleteExpense(ctx, int(e.ID))
		if err != nil {
			log.Fatalf("could not delete expense #%d", e.ID)
		}

		fmt.Printf("expense #%d deleted\n", e.ID)
	}

	// let's restore them
	fmt.Printf("restoring deleted expenses\n")
	for _, e := range newExpenses {
		err := client.RestoreExpense(ctx, int(e.ID))
		if err != nil {
			log.Fatalf("could not undelete expense #%d", e.ID)
		}

		fmt.Printf("expense #%d restored\n", e.ID)
	}

	// let's delete those test expenses again
	fmt.Printf("deleting restored expenses\n")
	for _, e := range newExpenses {
		err := client.DeleteExpense(ctx, int(e.ID))
		if err != nil {
			log.Fatalf("could not delete expense #%d", e.ID)
		}

		fmt.Printf("expense #%d deleted again\n", e.ID)
	}
}

func createExpenseBySharesExample(ctx context.Context, client splitwise.Client) {
	// let's do it inside of a group
	groups, err := client.GetGroups(ctx)
	if err != nil {
		log.Fatalf("could not get groups: %v", err)
	}

	if len(groups) == 1 {
		fmt.Printf("only one group (no group = 0), not doing anything\n")
		return
	}

	group := groups[0]
	for i := 1; i < len(groups) && group.ID == 0; i++ {
		group = groups[i]
	}

	fmt.Printf("selected group: #%d - %s\n", group.ID, group.Name)

	// let's see what users are inside the group
	users := group.Members

	// let's only do it for 2 users that are inside the group
	if len(users) < 3 {
		log.Printf("there are not enough users [%d users]\n", len(users))
		return
	}

	currentUser, err := client.GetCurrentUser(ctx)
	if err != nil {
		log.Fatalf("could not get current user: %v", err)
	}

	// all users but the current user
	for i, u := range users {
		if u.ID == currentUser.ID {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}

	users = users[:2]
	expUsers := []splitwise.ExpenseUser{}
	for _, user := range users {
		e := splitwise.ExpenseUser{
			Id:        user.ID,
			PaidShare: 0.0,
			OwedShare: 5000.0,
		}
		expUsers = append(expUsers, e)
	}

	expUsers = append(expUsers, splitwise.ExpenseUser{
		Id:        currentUser.ID,
		PaidShare: 10000.0,
		OwedShare: 0.0,
	})

	params := splitwise.CreateExpenseParams{
		splitwise.CreateExpenseRepeatInterval: "weekly",
	}

	expenses, err := client.CreateExpenseByShares(ctx, 10000, "should delete this", int(group.ID), params, expUsers)
	if err != nil {
		log.Fatalf("could not create expenses: %v", err)
	}

	fmt.Printf("%d expenses created\n", len(expenses))
	for _, e := range expenses {
		fmt.Printf("expense #%d - %s\n", e.ID, e.Description)
	}

	// let's try to update them
	fmt.Printf("updating created expenses\n")
	for _, e := range expenses {
		costValue, err := strconv.ParseFloat(e.Cost, 32)
		if err != nil {
			log.Fatalf("failed to convert cost to float: %v", err)
		}

		params = splitwise.CreateExpenseParams{
			splitwise.CreateExpenseRepeatInterval: "monthly",
		}

		updated, err := client.UpdateExpense(ctx, int(e.ID), costValue, e.Description, int(e.GroupId), params, nil)
		if err != nil {
			log.Fatalf("could not update expense #%d", e.ID)
		}

		fmt.Printf("expense #%d updated\n", e.ID)

		if len(updated) != 1 {
			log.Fatalf("the number of updated entries should be 1")
		}

		if updated[0].RepeatInterval != "monthly" {
			log.Fatalf("expense repeat interval was not updated to \"monthly\"")
		}
	}

	// let's delete those test expenses
	fmt.Printf("deleting created expenses\n")
	for _, e := range expenses {
		err := client.DeleteExpense(ctx, int(e.ID))
		if err != nil {
			log.Fatalf("could not delete expense #%d", e.ID)
		}

		fmt.Printf("expense #%d deleted\n", e.ID)
	}
}

func commentExamples(ctx context.Context, client splitwise.Client) {
	// let's do it inside of a group
	groups, err := client.GetGroups(ctx)
	if err != nil {
		log.Fatalf("could not get groups: %v", err)
	}

	if len(groups) == 1 {
		fmt.Printf("only one group (no group = 0), not doing anything\n")
		return
	}

	group := groups[0]
	for i := 1; i < len(groups) && group.ID == 0; i++ {
		group = groups[i]
	}

	fmt.Printf("selected group: #%d - %s\n", group.ID, group.Name)

	expenses, err := client.CreateExpenseEqualGroupSplit(ctx, 5000, "should delete - comments example", int(group.ID), nil)
	if err != nil {
		log.Fatalf("could not create expense for comments: %v", err)
	}

	defer func(ctx context.Context, client splitwise.Client, expenses []resources.Expense) {
		for _, e := range expenses {
			_ = client.DeleteExpense(ctx, int(e.ID))
		}
	}(ctx, client, expenses)

	if len(expenses) != 1 {
		log.Fatalf("created expenses should be only 1")
	}

	expense := expenses[0]

	commentId := createCommentsExample(ctx, client, int(expense.ID))
	defer func(ctx context.Context, client splitwise.Client, id int) {
		_, err := client.DeleteExpenseComment(ctx, id)
		if err != nil {
			log.Printf("could not delete comment #%d: %v\n", id, err)
		}
	}(ctx, client, commentId)

	queryCommentsExample(ctx, client, int(expense.ID))
}

func createCommentsExample(ctx context.Context, client splitwise.Client, expenseId int) int {
	comment, err := client.CreateExpenseComment(ctx, expenseId, "should delete this comment")
	if err != nil {
		log.Fatalf("could not create comment: %v", err)
	}

	return int(comment.ID)
}

func queryCommentsExample(ctx context.Context, client splitwise.Client, expenseId int) {
	comments, err := client.GetExpenseComments(ctx, expenseId)
	if err != nil {
		log.Fatalf("could not query comments: %v", err)
	}

	fmt.Printf("Comments\n")
	for _, c := range comments {
		fmt.Printf("Comment #%d: %s\n", c.ID, c.Content)
	}
}

func notificationsExamples(ctx context.Context, client splitwise.Client) {
	ns, err := client.GetNotifications(ctx, nil)
	if err != nil {
		log.Fatalf("could not get notifications: %v", err)
	}

	fmt.Printf("Notifications\n")
	for _, n := range ns {
		fmt.Printf("Notification #%d: %s\n", n.ID, n.Content)
	}
}

func currenciesExamples(ctx context.Context, client splitwise.Client) {
	currs, err := client.GetCurrencies(ctx)
	if err != nil {
		log.Fatalf("could not get currencies: %v", err)
	}

	fmt.Printf("Currencies\n")
	for _, c := range currs {
		fmt.Printf("Currency %s - %s\n", c.CurrencyCode, c.Unit)
	}
}

func categoriesExamples(ctx context.Context, client splitwise.Client) {
	cats, err := client.GetCategories(ctx)
	if err != nil {
		log.Fatalf("could not get categories: %v", err)
	}

	fmt.Printf("Categories\n")
	for _, c := range cats {
		fmt.Printf("Category    %d - %s\n", c.ID, c.Name)
		for _, sc := range c.Subcategories {
			fmt.Printf("Subcategory %d - %s\n", sc.ID, sc.Name)
		}
	}
}

func parsedExpenseExamples(ctx context.Context, client splitwise.Client) {
	friends, err := client.GetFriends(ctx)
	if err != nil {
		log.Fatalf("could not get friends: %v", err)
	}

	if len(friends) == 0 {
		log.Printf("no friends to create the expense :(")
		return
	}

	friend := friends[0]

	input := fmt.Sprintf("I owe %s %s $5000 COP", friend.FirstName, friend.LastName)
	fmt.Printf("Input: %s\n", input)
	pexp, err := client.ParseSentenceIntoExpenseFreeform(ctx, input, false)
	if err != nil {
		log.Fatalf("could not parse sentense into expense: %v", err)
	}

	defer func() {
		_ = client.DeleteExpense(ctx, int(pexp.Expense.ID))
	}()

	var users []string
	for _, u := range pexp.Expense.Users {
		users = append(users, u.User.FirstName)
	}

	fmt.Printf("Parsed expense #%d [valid: %v] - $%s - %s - %v\n",
		pexp.Expense.ID, pexp.Valid,
		pexp.Expense.Cost, pexp.Expense.Description,
		users,
	)
}
