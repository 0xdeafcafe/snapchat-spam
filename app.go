package main

import (
	"./snapchat"
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
)

func GetArguments() (string, string, string, int, error) {
	flag.Parse()

	token := flag.Arg(0)
	usernameToSend := flag.Arg(1)
	usernameToAbuse := flag.Arg(2)
	count := flag.Arg(3)

	counti, err := strconv.Atoi(count)
	if err != nil {
		return "", "", "", 0, errors.New("The number of times to spam someone needs to be a number..")
	}
	if token == "" {
		return "", "", "", 0, errors.New("You need a token, son..")
	}
	if usernameToSend == "" {
		return "", "", "", 0, errors.New("You need a sender, dude..")
	}
	if usernameToAbuse == "" {
		return "", "", "", 0, errors.New("You need a victim, dude..")
	}
	if counti == 0 {
		return "", "", "", 0, errors.New("You need to specify the number of times to spam the vicim..")
	}

	return token, usernameToSend, usernameToAbuse, counti, nil
}

func main() {
	token, usernameToSend, usernameToAbuse, count, err := GetArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := snapchat.Prep()
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < count; i++ {
		mediaId, _, err := snapchat.UploadMedia(token, data, usernameToSend)
		if err != nil {
			log.Fatal(err)
		}
		_, err = snapchat.SendMedia(token, usernameToAbuse, usernameToSend, mediaId)
		if err != nil {
			fmt.Println(fmt.Sprintf("[%d] - Error Sending Snap %s", i, usernameToAbuse))
		} else {
			fmt.Println(fmt.Sprintf("[%d] - Snap Sent to %s", i, usernameToAbuse))
		}
	}
}
