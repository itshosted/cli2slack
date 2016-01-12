package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pelletier/go-toml"
	"net/http"
	"net/url"
	"os"
	"time"
	"io/ioutil"
	"errors"
)

type WebhookAttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}
type WebhookAttachment struct {
	Fallback string                   `json:"fallback"`
	Pretext  string                   `json:"pretext"`
	Text     string                   `json:"text"`
	Fields   []WebhookAttachmentField `json:"fields"`
}
type Webhook struct {
	Channel     string              `json:"channel"`
	Username    string              `json:"username"`
	IconEmoji   string              `json:"icon_emoji"`
	Text        string              `json:"text"`
	Attachments []WebhookAttachment `json:"attachments"`
}

func main() {

	tomlfile := ""
	message := ""
	channel := ""
	flag.StringVar(&tomlfile, "config", "cli2slack.conf", "Location of the config file")
	flag.StringVar(&channel, "c", "testen", "Channel to send message to")
	flag.StringVar(&message, "m", "No Message :)", "Message to send to slack channel")
	flag.Parse()

	config, err := toml.LoadFile(tomlfile)
	if err != nil {
		fmt.Println("Error ", err.Error())
	} else {
		// retrieve data directly
		config_Url := config.Get("url").(string)
		config_username := config.Get("username").(string)
		config_iconEmoji := config.Get("iconEmoji").(string)

		//fmt.Println("Url is ", Url, ". Username is ", username, ". Channel is ", channel)

		h, e := os.Hostname()
		if e != nil {
			panic(e)
		}
		
		// Create JSON
		str, e := json.Marshal(Webhook{
			Text:      "",
			Channel:   "#" + channel,
			Username:  config_username,
			IconEmoji: config_iconEmoji,
			Attachments: []WebhookAttachment{WebhookAttachment{
				Fallback: "",
				Pretext:  "",
				Text:     message,
				Fields: []WebhookAttachmentField{WebhookAttachmentField{
					Title: "Hostname",
					Value: h,
					Short: true,
				}, WebhookAttachmentField{
					Title: "Date",
					Value: time.Now().Format("2006-Jan-02 15:04"),
					Short: true,
				}},
			}},
		})
		
		// Post JSON to slack
		res, e := http.PostForm(
			config_Url, url.Values{"payload": {string(str)}},
		)
		if e != nil {
			panic(e)
		}

		if res.StatusCode != 200 {
			defer res.Body.Close()
			txt, e := ioutil.ReadAll(res.Body)
			if e != nil {
				panic(e)
			}
			fmt.Println( errors.New(fmt.Sprintf("HTTP=%d, txt=%s", res.StatusCode, string(txt))))
		}
	}

}
