package main

import (
	"bytes"
	"configrd"
	"encoding/json"
	"logger"
	"net/http"
	"os"
)

type AliCallback struct {
	push_data  map[string]string
	repository map[string]string
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
		push_data:  make(map[string]string),
		repository: make(map[string]string),
	}
	jdecoder := json.NewDecoder(r.Body)
	err := jdecoder.Decode(&cont)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to unmarshal callback,", err)
		return
	}
	noti := Notification{
		Heading: "构建完毕",
		Content: "",
	}
	noti.Content = cont.repository["repo_full_name"] + ":" + cont.push_data["tag"] + " # " + cont.push_data["pushed_at"]
	marshed, _ := json.Marshal(noti)
	_, httperr := http.Post(conf["CollectorAddr"], "application/json", bytes.NewReader(marshed))
	if httperr != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to post to collector,", httperr)
	}
}
