package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/site", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}
		request, err := http.NewRequest("GET", "http://"+r.FormValue("host"), nil)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(request)
		fmt.Println(resp)
		resp.Body.Close()
		fmt.Fprintf(w, "foda-se")
	})
	http.ListenAndServe(":7777", nil)
}
