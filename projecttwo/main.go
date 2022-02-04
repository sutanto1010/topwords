package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	//filePath := `../../training-parallel/training/big-one.txt`
	//filePath := `../../training-parallel/training/europarl-v6.fr-en.en`
	//filePath := `./hello.txt`
	filePath := `./GoLang_Test.txt`

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	prevTextTail := ""
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp := &Response{}
	stop := false
	for !stop {
		usedText := ""
		buff := make([]byte, 1024*100)
		_, err := file.Read(buff)
		if err != nil {
			usedText = prevTextTail + string(buff)
			resp, _ = postTextToServer(client, usedText, resp.SessionID)
			stop = true
		}
		text := string(buff)
		lextLength := len(text)
		lastSpaceIndex := strings.LastIndex(text, " ")

		if lastSpaceIndex != -1 && lastSpaceIndex < lextLength-1 {
			usedText = prevTextTail + text[:lastSpaceIndex]
			prevTextTail = text[lastSpaceIndex+1:]
		} else if lastSpaceIndex == -1 {
			prevTextTail = prevTextTail + text
			usedText = ""
		} else {
			usedText = prevTextTail + text
			prevTextTail = ""
		}
		resp, _ = postTextToServer(client, usedText, resp.SessionID)
	}
	//fmt.Println("resp: ", resp)
	for _, line := range resp.TopWords {
		fmt.Println(line.Word, " : ", line.Count)
	}
}

func postTextToServer(client *http.Client, text, sessionID string) (*Response, error) {
	url := "http://localhost:8080/topwords"
	payloadMap := map[string]interface{}{
		"text":       text,
		"session_id": sessionID,
	}
	payload, _ := json.Marshal(payloadMap)
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res Response
	err = json.Unmarshal(body, &res)
	return &res, err
}

type TopWord struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type Response struct {
	TopWords  []TopWord `json:"top_words"`
	SessionID string    `json:"session_id"`
}
