package main

import (
	"bytes"
	"configrd"
	"encoding/json"
	"fmt"
	"logger"
	"net/http"
	"os"
)

type AliCallback struct {
	Push_data  map[string]string
	Repository map[string]string
}

type Notification struct {
	Heading string
	Content string
}

var conf map[string]string

func main() {
	var confile = configrd.Config(os.Args[1])
	conf = confile.ReadConfig()
	_, mpok := conf["ListenPort"]
	if !mpok {
		panic("ListenPort has to be exist in config map!")
	}
	_, mpok = conf["CollectorAddr"]
	if !mpok {
		panic("CollectorAddr has to be exist in config map!")
	}

	http.HandleFunc("/ccallback", onCallback)
	logger.Log.Logln(logger.LEVEL_PANIC, "Listen", http.ListenAndServe(":"+conf["ListenPort"], nil))
}

func onCallback(w http.ResponseWriter, r *http.Request) {
	cont := AliCallback{
		Push_data:  make(map[string]string),
		Repository: make(map[string]string),
	}
	jdecoder := json.NewDecoder(r.Body)
	err := jdecoder.Decode(&cont)
	fmt.Println("Recv ", cont)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to unmarshal callback,", err)
		return
	}
	noti := Notification{
		Heading: "构建完毕",
		Content: "",
	}
	noti.Content = cont.Repository["repo_full_name"] + ":" + cont.Push_data["tag"] + " # " + cont.Push_data["pushed_at"]
	marshed, _ := json.Marshal(noti)
	_, httperr := http.Post(conf["CollectorAddr"], "application/json", bytes.NewReader(marshed))
	if httperr != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to post to collector,", httperr)
	}
}
