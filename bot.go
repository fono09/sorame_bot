package main

import (
	"errors"
	"fmt"
	"config"
	"model"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var me twitter.User
var client *twitter.Client

func sorameCollector(tweet *twitter.Tweet) {
	if tweet.User.ID == me.ID || strings.Index(tweet.Text, "RT") != -1 {
		return
	}

	log.Printf("Tweet: @%s %s\n", tweet.User.ScreenName, tweet.Text)
	sorame, err := model.NewSorameFromTweet(tweet)
	if err == nil {
		log.Println("=========================")
		log.Printf("New Sorame: %#v\n", sorame)
		log.Println("=========================")
		err = sorame.Save()
		status := ""
		if err != nil {
			status = fmt.Sprintf("@%s 空目の登録に失敗しました。そもそも空目ではない可能性があります。", tweet.User.ScreenName)
		} else {
			status = fmt.Sprintf(". @%s さんが %s を %s に空目しました。削除をご希望の場合、sorame_bot_deleteと返信してください。", tweet.User.ScreenName, sorame.Before, sorame.After)
		}
		client.Statuses.Update(status, &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID})
	}

}

func findTweetByID(id int64) *twitter.Tweet {
	tweets, _, err := client.Statuses.Lookup([]int64{id}, nil)
	if err != nil || len(tweets) == 0 {
		return nil
	}
	log.Printf("findTweetById(%d) len(tweets) :%d", id, len(tweets))
	return &tweets[0]
}

func getTweetInReplyTo(tweet *twitter.Tweet) *twitter.Tweet {
	inReplyToStatusID := tweet.InReplyToStatusID
	return findTweetByID(inReplyToStatusID)
}

func sorameRemover(tweet *twitter.Tweet) {
	if strings.Index(tweet.Text, "sorame_bot_delete") == -1 {
		return
	} else if tweet.User.ID == me.ID {
		return
	}

	deleted := false

	err := errors.New("")
	beforeTweet := getTweetInReplyTo(tweet)
	log.Printf("beforeTweet(0): %#v\n", beforeTweet)
	if beforeTweet != nil && beforeTweet.User.ID == me.ID {
		beforeTweet = getTweetInReplyTo(beforeTweet)
		log.Printf("beforeTweet(1): %#v\n", beforeTweet)
		if beforeTweet != nil && beforeTweet.User.ID != me.ID && beforeTweet.User.ID == tweet.User.ID {
			sorame := &model.Sorame{ID: beforeTweet.ID}
			err = sorame.RemoveByID()
			if err == nil {
				deleted = true
			}
		}
	}

	status := ""
	if deleted {
		status = fmt.Sprintf("@%s 削除しました", tweet.User.ScreenName)
	} else {
		status = fmt.Sprintf("@%s 削除に失敗。%s", tweet.User.ScreenName, err.Error())
	}
	client.Statuses.Update(status, &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID})
}

var tlTweetCount = 0

func getUserFromID(id int64) *twitter.User {
	ulp := &twitter.UserLookupParams{UserID: []int64{id}}
	users, _, err := client.Users.Lookup(ulp)
	if err != nil {
		return nil
	}
	return &users[0]
}
func soramePickupper() {
	tlTweetCount++
	if tlTweetCount/(2*me.FollowersCount) > 1 {
		//if tlTweetCount == 5 {
		sorame := model.Sorame{}
		err := sorame.RandomGet()
		if err == nil {
			status := fmt.Sprintf("%s を %s に空目 @%s", sorame.Before, sorame.After, getUserFromID(sorame.UserID).ScreenName)
			client.Statuses.Update(status, nil)
		}
		tlTweetCount = 0
	}

}

func tweetListener(tweet *twitter.Tweet) {
	go sorameCollector(tweet)
	go sorameRemover(tweet)
	go soramePickupper()
}

func followResponder(event *twitter.Event) {
	if event.Event == "follow" || event.Target.ID == me.ID {
		follow := true
		frsParams := twitter.FriendshipCreateParams{UserID: event.Source.ID, Follow: &follow}
		_, _, err := client.Friendships.Create(&frsParams)
		if err != nil {
			client.Statuses.Update(fmt.Sprintf("@%s フォロー返しました。", event.Source.ScreenName), nil)
		} else {
			client.Statuses.Update(fmt.Sprintf("@%s フォロー返しに失敗しました。", event.Source.ScreenName), nil)
		}
	}
}

func eventListener(event *twitter.Event) {
	followResponder(event)
}

func initTwitter() {
	conf := oauth1.NewConfig(config.Config.App.Token.ConsumerKey, config.Config.App.Token.ConsumerSecret)
	token := oauth1.NewToken(config.Config.App.Token.AccessToken, config.Config.App.Token.AccessTokenSecret)
	httpClient := conf.Client(oauth1.NoContext, token)

	client = twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	acc, _, _ := client.Accounts.VerifyCredentials(verifyParams)
	me = *acc
}
func initStream() {
	params := &twitter.StreamUserParams{
		With:          "followings",
		StallWarnings: twitter.Bool(true),
	}

	// Demultiplexer
	demux := twitter.NewSwitchDemux()
	demux.Tweet = tweetListener
	demux.Event = eventListener
	demux.DM = func(dm *twitter.DirectMessage) {
		log.Printf("DM: %#v", dm)
	}
	stream, err := client.Streams.User(params)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	demux.HandleChan(stream.Messages)
}

func captureSignal() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	for s := range signalCh {
		switch s {
		case syscall.SIGINT:
			status := fmt.Sprintf("動作停止しました %s", time.Now())
			client.Statuses.Update(status, nil)
			log.Println("Server will stop...")
			os.Exit(0)
		}
	}
}

func main() {
	go captureSignal()
	model.InitDB()
	initTwitter()
	status := fmt.Sprintf("動作開始しました %s", time.Now())
	client.Statuses.Update(status, nil)
	initStream()
}
