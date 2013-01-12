package main

import (
  //"fmt"
  "encoding/json"
  "log"
  "net/http"
  "os"
  "text/template"
  "github.com/gorilla/mux"
)

func ViewHandler ( w http.ResponseWriter, r *http.Request ) {
  //fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
  vars := mux.Vars(r)
  id := vars["id"]

  var page Page

  file, err := os.Open("/home/stephen/hacking/server/data/" + id) // For read access.
  if err != nil {
    log.Fatal(err)
  }

  enc := json.NewDecoder(file)
  enc.Decode(&page)

  tmpl, err := template.ParseFiles("/home/stephen/hacking/server/templates/layout.html")
  
  if err != nil { panic(err) }
  err = tmpl.Execute(w, page)
  if err != nil { panic(err) }

}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/view/{id:[a-z0-9-]+}", ViewHandler)

    http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}
