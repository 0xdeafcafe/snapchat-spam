package snapchat

import (
	"../crypto"
	"bytes"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Token string

const (
	pattern           = "0001110111101110001111010101111011010001001110011000110001000110"
	secret            = "iEk21fuwZApXlz93750dmW22pw389dPwOk"
	staticToken       = "m198sOkJEn37DjqZ32lpRu76xmw288xSQ9"
	userAgent         = "Snapchat/8.1.0.11 (iPhone5,2; iOS 8.1; gzip)"
	blobEncryptionKey = "M02cnQ51Ji97vwT4"
)

func Prep(filePath string) ([]byte, error) {
	imageData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	data := crypto.Encrypt(imageData, blobEncryptionKey)
	return data, nil
}

func (t Token) UploadMedia(data []byte, username string) (string, bool, error) {
	mediaId := generateMediaId(username)

	params := map[string]string{
		"media_id": mediaId,
		"type":     "0",
		"zipped":   "0",
	}
	resp, err := sendMultipartPostRequest("bq/upload", username, t, params, data)
	if err != nil {
		log.Fatal(err)
	}
	return mediaId, (resp.StatusCode == http.StatusOK), nil
}

func (t Token) SendMedia(recipient string, username string, mediaId string) (bool, error) {
	params := map[string]string{
		"country_code": "GB",
		"media_id":     mediaId,
		"recipients":   fmt.Sprintf("[%s]", recipient),
		"reply":        "0",
		"time":         "10",
		"type":         "0",
		"zipped":       "0",
	}
	resp, err := sendPostRequest("loq/send", username, t, params)
	if err != nil {
		log.Fatal(err)
	}
	// bdy, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(bdy))
	return (resp.StatusCode == http.StatusOK), nil
}

func (t Token) SendChatMedia(data []byte, username string, recipient string) (bool, error) {
	id := generateId()

	params := map[string]string{
		"conversation_id": fmt.Sprintf("%s~%s", recipient, username),
		"id":              id,
		"recipient":       recipient,
		"type":            "IMAGE",
	}
	resp, err := sendMultipartPostRequest("bq/upload_chat_media", username, t, params, data)
	if err != nil {
		log.Fatal(err)
	}
	return (resp.StatusCode == http.StatusOK), nil
}

func generateRequestToken(token Token, timestamp string) string {
	s1 := fmt.Sprintf("%s%s", secret, token)
	s2 := fmt.Sprintf("%s%s", timestamp, secret)

	s3 := crypto.Sha256(s1)
	s4 := crypto.Sha256(s2)

	output := ""
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '0' {
			output += string(s3[i])
		} else {
			output += string(s4[i])
		}
	}
	return output
}

func sendMultipartPostRequest(endpoint string, username string, token Token, params map[string]string, data []byte) (*http.Response, error) {
	timestamp := generateTimestamp()
	requestToken := generateRequestToken(token, timestamp)

	// Set our params aswell
	params["req_token"] = requestToken
	params["timestamp"] = timestamp
	params["username"] = username

	// Create the multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range params {
		writer.WriteField(key, value)
	}

	// Are we adding a file to upload?
	if data != nil {
		part, err := writer.CreateFormFile("data", "data")
		if err != nil {
			return nil, err
		}
		part.Write(data)
	}
	writer.Close()

	// Create the request
	req, err := http.NewRequest("POST", fmt.Sprintf("https://feelinsonice-hrd.appspot.com/%s", endpoint), &body)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "en;q=1")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Content-Length", string(body.Len()))

	// Do the request, return response
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func sendPostRequest(endpoint string, username string, token Token, params map[string]string) (*http.Response, error) {
	timestamp := generateTimestamp()
	requestToken := generateRequestToken(token, timestamp)

	// Set our params aswell
	params["req_token"] = requestToken
	params["timestamp"] = timestamp
	params["username"] = username

	// Create the multipart form
	data := url.Values{}
	for key, value := range params {
		data.Add(key, value)
	}

	// Create the request
	req, err := http.NewRequest("POST", fmt.Sprintf("https://feelinsonice-hrd.appspot.com/%s", endpoint), bytes.NewBufferString(data.Encode()))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "en;q=1")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	// Do the request, return response
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func generateTimestamp() string {
	timestampi := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%d", timestampi)
}

func generateMediaId(username string) string {
	id, _ := uuid.NewV4()
	return strings.ToUpper(username + "~" + id.String())
}

func generateId() string {
	id, _ := uuid.NewV4()
	return strings.ToUpper(id.String())
}
