package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Conf struct {
	Port int    `json:"port,omitempty"`
	Resp string `json:"resp,omitempty"`
}

type CustomResponse struct {
	Document string `json:"document"`
}

var (
	conf Conf
	port int = 8001
)

func readConf(c *Conf) error {
	b, err := ioutil.ReadFile("conf.json")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, c); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := readConf(&conf); err != nil {
		log.Fatalf("cannot read conf.json, err: %v", err)
	}

	http.HandleFunc("/api", HelloServer)
	http.HandleFunc("/update", UpdateConfServer)
	fmt.Println("running on ", port)

	// TODO: add graceful shutdown
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(CustomResponse{Document: conf.Resp})
}

func UpdateConfServer(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "Only POST requests allowed.", http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("before updating conf: %#v\n", conf)
	if err := json.NewDecoder(r.Body).Decode(&conf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("after updating conf: %#v\n", conf)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(CustomResponse{Document: "Config updated."})
}
