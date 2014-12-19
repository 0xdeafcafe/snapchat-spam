package main

import (
	"./snapchat"
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
)

const (
	// Image path of the image of Taylor Swift (exclusive) we're using to spam.
	BaePath = "./baelor'd.jpg"
)

func getArguments() (string, string, string, int, error) {
	flag.Parse()

	token := flag.Arg(0)
	usernameToSend := flag.Arg(1)
	usernameToAbuse := flag.Arg(2)
	count := flag.Arg(3)

	counti, err := strconv.Atoi(count)
	if err != nil {
		return "", "", "", 0, errors.New("the number of times to spam someone needs to be a number")
	}
	if token == "" {
		return "", "", "", 0, errors.New("you need a token, son")
	}
	if usernameToSend == "" {
		return "", "", "", 0, errors.New("you need a sender, dude")
	}
	if usernameToAbuse == "" {
		return "", "", "", 0, errors.New("you need a victim, dude")
	}
	if counti == 0 {
		return "", "", "", 0, errors.New("you need to specify the number of times to spam the vicim")
	}

	return token, usernameToSend, usernameToAbuse, counti, nil
}

func main() {
	tokenString, usernameToSend, usernameToAbuse, count, err := GetArguments()
	if err != nil {
		log.Fatal(err)
	}

	data, err := snapchat.Prep(BaePath)
	if err != nil {
		log.Fatal(err)
	}

	token := snapchat.Token(tokenString)
	for i := 0; i < count; i++ {
		mediaID, _, err := token.UploadMedia(data, usernameToSend)
		if err != nil {
			log.Fatal(err)
		}

		response, err := token.SendMedia(usernameToAbuse, usernameToSend, mediaID)
		if err != nil || !response {
			fmt.Println(fmt.Sprintf("[%d] - Error Sending Snap to %s", i, usernameToAbuse))
			log.Fatal(err)
		} else {
			fmt.Println(fmt.Sprintf("[%d] - Snap Sent to %s", i, usernameToAbuse))
		}
	}
}
