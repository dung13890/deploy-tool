package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const key = "r3kqouWSZ1"

type params struct {
	Project string `json:"project"`
	Value   string `json:"value"`
	Result  bool   `json:"result"`
}

func main() {
	http.HandleFunc("/openssl", opensslHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	fmt.Printf("Starting server at port 1308\n")
	if err := http.ListenAndServe(":1308", nil); err != nil {
		log.Fatal(err)
	}
}

func opensslHandler(w http.ResponseWriter, r *http.Request) {
	k := r.Header.Get("X-SecretKey")
	if k != key {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	d := json.NewDecoder(r.Body)
	p := &params{}
	err := d.Decode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := p.writeFile(); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", p)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func (p *params) writeFile() error {
	fi, err := os.OpenFile("public/output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	defer fi.Close()
	if err != nil {
		return err
	}
	rs := "Failed"
	if p.Result == true {
		rs = "Pass"
	}
	fi.WriteString(fmt.Sprintf("%s,%s,%s\n", p.Project, p.Value, rs))

	return nil
}
