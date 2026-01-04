package commands

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"strconv"
	"time"

	"github.com/Rachit-Gandhi/gator/internal/database"
	"github.com/Rachit-Gandhi/gator/internal/rss"
	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	_, ok := s.Db.GetUser(context.Background(), cmd.StringArgs[0])
	if ok != nil {
		return fmt.Errorf("user doesn't exist in the database")
	}
	err := s.Cfg.SetUser(cmd.StringArgs[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set as: %v\n", cmd.StringArgs[0])
	return nil
}

func Register(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	_, ok := s.Db.GetUser(context.Background(), cmd.StringArgs[0])
	if ok == nil {
		return fmt.Errorf("user with that username already exists")
	}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.StringArgs[0],
	}
	_, err := s.Db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("user creation in db failed, %w", err)
	}
	err = s.Cfg.SetUser(cmd.StringArgs[0])
	if err != nil {
		return err
	}
	fmt.Printf("new user has been created as: %v\n", cmd.StringArgs[0])
	return nil
}

func Reset(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 0 {
		return fmt.Errorf("the reset handler doesn't accept cli argument")
	}
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting all users: %w\n", err)
	}
	fmt.Println("all users are reset")
	return nil
}

func GetUser(s *State, cmd Command) (database.User, error) {
	if len(cmd.StringArgs) != 1 {
		return database.User{}, fmt.Errorf("the user handler needs one argument username of type string")
	}
	user, err := s.Db.GetUser(context.Background(), cmd.StringArgs[0])
	if err != nil {
		return database.User{}, fmt.Errorf("error getting user with this username: %w", err)
	}
	return user, nil
}

func GetUsers(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 0 {
		return fmt.Errorf("the users handler doesn't accept cli argument")
	}
	usrs, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting all users: %w\n", err)
	}
	currentUser := s.Cfg.GetUser()
	for _, usr := range usrs {
		fmt.Print("* ")
		if usr.Name != currentUser {
			fmt.Printf("%v\n", usr.Name)
		} else {
			fmt.Printf("%v (current)\n", currentUser)
		}
	}
	return nil
}

func Aggregate(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 1 {
		return fmt.Errorf("the users handler accepts one argument time_between_reqs")
	}
	time_between_reqs, err := time.ParseDuration(cmd.StringArgs[0])
	if err != nil {
		return fmt.Errorf("error parsing time betweeen reqs: %w", err)
	}
	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("error scraping feed: %v", err)
			continue
		}
	}
}

func AddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.StringArgs) != 2 {
		return fmt.Errorf(("add feed handlers expects 2 arguments: name & url"))
	}
	id := uuid.New()
	user_id := user.ID
	name, url := cmd.StringArgs[0], cmd.StringArgs[1]
	newFeed := database.AddFeedParams{
		ID:     id,
		UserID: user_id,
		Name:   name,
		Url:    url,
	}
	_, err := s.Db.AddFeed(context.Background(), newFeed)
	if err != nil {
		return fmt.Errorf("error adding feed: %w", err)
	}
	feed_follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user_id,
		FeedID:    newFeed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), feed_follow)
	if err != nil {
		return fmt.Errorf("error adding feed follow post adding feed: %w", err)
	}
	fmt.Printf("feed added: %v\n", newFeed)
	return nil
}

func GetFeeds(s *State, cmd Command) error {
	if len(cmd.StringArgs) != 0 {
		return fmt.Errorf("get feeds doesn't accept any arguments")
	}
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}
	for _, feed := range feeds {
		userName, err := s.Db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting username: %w", err)
		}
		fmt.Printf("Name: %v, URL: %v, userName: %v\n", feed.Name, feed.Url, userName.Name)
	}
	return nil
}

func CreateFeedFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.StringArgs) != 1 {
		return fmt.Errorf("create feed expects one argument, feed_url")
	}
	feed_url := cmd.StringArgs[0]
	feed, err := s.Db.GetFeedByUrl(context.Background(), feed_url)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no feed entry found to follow: %w", err)
	} else if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}
	feed_follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), feed_follow)
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}
	fmt.Printf("feed follow created between %v user and %v feed", user.Name, feed.Name)
	return nil
}

func GetFeedFollowsForUser(s *State, cmd Command, user database.User) error {
	if len(cmd.StringArgs) > 0 {
		return fmt.Errorf("following doesn't accept any args.")
	}
	feedsFollowed, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feeds followed by user: %w", err)
	}
	fmt.Printf("* %v\n", s.Cfg.GetUser())
	for _, feed := range feedsFollowed {
		feed_name, err := s.Db.GetFeedNameById(context.Background(), feed.FeedID)
		if err != nil {
			return fmt.Errorf("error getting feed name: %w", err)
		}
		fmt.Printf("    - %v\n", feed_name)
	}
	return nil
}

func DeleteFeedFollowsPair(s *State, cmd Command, user database.User) error {
	if len(cmd.StringArgs) != 1 {
		return fmt.Errorf("unfollow expects one argument, feedurl")
	}
	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.StringArgs[0])
	if err != nil {
		return fmt.Errorf("error getting feed from url: %w", err)
	}
	deletePair := database.DeleteFeedFollowsPairParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.Db.DeleteFeedFollowsPair(context.Background(), deletePair)
	if err != nil {
		return fmt.Errorf("error deleting the pair: %w", err)
	}
	return nil
}

func BrowseFeeds(s *State, cmd Command, user database.User) error {
	if len(cmd.StringArgs) > 1 {
		return fmt.Errorf("browse expects max one argument limit")
	}
	limit := 2
	if len(cmd.StringArgs) == 1 {
		parsed, err := strconv.Atoi(cmd.StringArgs[0])
		if err != nil {
			return fmt.Errorf("limit must be a number: %w", err)
		}
		limit = parsed
	}
	userAndLimit := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.Db.GetPostsForUser(context.Background(), userAndLimit)
	if err != nil {
		return fmt.Errorf("error getting posts for the user: %w", err)
	}
	for i, post := range posts {
		fmt.Printf("Post: %v\n\tTitle:%v\n\tPostURL:%v\n\tDescription:%v\n", i, post.Title, post.PostUrl, post.PostDescription)
	}
	return nil
}

func limitWord(word string, limit int) string {
	runes := []rune(word)
	if len(runes) > limit {
		return string(runes[:limit])
	}
	return word
}

func parseRSSDate(dateStr string) (time.Time, error) {
	formats := []string{
		"Mon, 02 Jan 2006 15:04:05 -0700", // RFC1123Z
		"Mon, 02 Jan 2006 15:04:05 MST",   // RFC1123
		time.RFC822,
		time.RFC822Z,
		time.RFC3339,
	}
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func scrapeFeeds(s *State) error {
	feed, err := s.Db.GetNextFeedtoFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %w", err)
	}
	r, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching the feed: %w", err)
	}
	fmt.Printf("Scraping feed: %s (ID: %v)\n", feed.Name, feed.ID)
	for _, post := range r.Channel.Item {
		fmt.Printf("fetching: %v\n", post.Title)
		t, err := parseRSSDate(post.PubDate)
		if err != nil {
			fmt.Printf("error parsing time resorting to default current time: %v\n", err)
			t = time.Now()
		}

		newPost := database.CreatePostParams{
			ID:              uuid.New(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			Title:           limitWord(html.UnescapeString(post.Title), 100),
			PostUrl:         limitWord(html.UnescapeString(post.Link), 500),
			PostDescription: limitWord(html.UnescapeString(post.Description), 500),
			PublishedAt:     t,
			FeedID:          feed.ID,
		}
		_, err = s.Db.CreatePost(context.Background(), newPost)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return fmt.Errorf("error adding post: %w", err)
		}
	}
	err = s.Db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("error marking the feed fetched: %w", err)
	}
	return nil
}
