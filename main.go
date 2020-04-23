package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

func main() {
	var (
		tl   = flag.String("tl", "zh-TW", "Google translate languages")
		text = flag.String("text", "", "The text you want to play")
	)

	flag.Parse()

	if *text == "" {
		fmt.Println("Please input text.")
		return
	}

	encodeText := url.QueryEscape(*text)
	endpoint := fmt.Sprintf("https://translate.google.com.vn/translate_tts?ie=UTF-8&q=%s&tl=%s&client=tw-ob", encodeText, *tl)
	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	dec, data, _ := minimp3.DecodeFull(body)
	context, err := oto.NewContext(dec.SampleRate, dec.Channels, 2, 1<<11)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer context.Close()
	player := context.NewPlayer()
	defer player.Close()
	if _, err := player.Write(data); err != nil {
		fmt.Println(err)
		return
	}
}
