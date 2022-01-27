package main

import (
	"encoding/json"
	"fmt"
	"github.com/buraksekili/bak"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Conf struct {
	Port int    `json:"port"`
	Resp string `json:"resp"`
}

type CustomResponse struct {
	Document string `json:"document"`
}

var (
	conf Conf = Conf{}
	port int  = 8001
)

func main() {
	w := bak.New(bak.Conf{File: "conf.json", Duration: 1 * time.Second})
	go func() {
		ch := w.Watch()
		for range ch {
			fmt.Println("changes detected in conf.json file.")
			if err := readConf(); err != nil {
				log.Fatalf("cannot read conf.json, err: %v", err)
			}
		}
	}()
	if err := readConf(); err != nil {
		log.Fatalf("cannot read conf.json, err: %v", err)
	}

	http.HandleFunc("/api", HelloServer)
	fmt.Println("running on ", port)

	// TODO: add graceful shutdown
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func readConf() error {
	b, err := ioutil.ReadFile("conf.json")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &conf); err != nil {
		return err
	}
	return nil
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	w.Header().Set("Content-Type", "application/json")
	resp := conf.Resp
	if resp == "" {
		resp = "default response"
	}
	res := CustomResponse{Document: resp}
	json.NewEncoder(w).Encode(res)
}
