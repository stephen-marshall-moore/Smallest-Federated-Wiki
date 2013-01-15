package main

import (
  //"fmt"
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "text/template"
  "github.com/gorilla/mux"
)

func StaticHandler ( w http.ResponseWriter, r *http.Request ) {
  vars := mux.Vars(r)
  fname := "/home/stephen/hacking/fedwiki/client/" + vars["fn"]
  log.Println( "StaticHandler: " + fname )

  http.ServeFile(w,r,fname)
}

func JsHandler ( w http.ResponseWriter, r *http.Request ) {
  vars := mux.Vars(r)
  fname := "/home/stephen/hacking/fedwiki/client/js/" + vars["fn"]
  log.Println( "JsHandler: " + fname )

  http.ServeFile(w,r,fname)
}

func JsSubHandler ( w http.ResponseWriter, r *http.Request ) {
  vars := mux.Vars(r)
  fname := "/home/stephen/hacking/fedwiki/client/js/" + vars["sub"] + "/" + vars["fn"]
  log.Println( "JsSubHandler: " + fname )

  http.ServeFile(w,r,fname)
}


func PluginHandler ( w http.ResponseWriter, r *http.Request ) {
  vars := mux.Vars(r)
  fname := "/home/stephen/hacking/fedwiki/client/plugins/" + vars["fn"]
  log.Println( "PluginHandler: " + fname )

  http.ServeFile(w,r,fname)
}

func PluginSubHandler ( w http.ResponseWriter, r *http.Request ) {
  vars := mux.Vars(r)
  fname := "/home/stephen/hacking/fedwiki/client/plugins/" + vars["sub"] + "/" + vars["fn"]
  log.Println( "PluginSubHandler: " + fname )

  http.ServeFile(w,r,fname)
}


func JsonHandler ( w http.ResponseWriter, r *http.Request ) {
  vars := mux.Vars(r)
  slug := "/home/stephen/hacking/fedwiki/server/data/" + vars["slug"]
  log.Println( "JsonHandler: " + slug )

  http.ServeFile(w,r,slug)
}

func ViewHandler ( w http.ResponseWriter, r *http.Request ) {
  //fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
  vars := mux.Vars(r)
  id := vars["id"]

  var page Page

  file, err := os.Open("/home/stephen/hacking/fedwiki/server/data/" + id) // For read access.
  if err != nil {
    log.Fatal(err)
  }

  enc := json.NewDecoder(file)
  enc.Decode(&page)

  tmpl, err := template.ParseFiles("/home/stephen/hacking/fedwiki/server/templates/layout.html")
  
  if err != nil { panic(err) }
  err = tmpl.Execute(w, page)
  if err != nil { panic(err) }

}

func MapHandler ( w http.ResponseWriter, r *http.Request ) {
  w.Header().Set( "Content-Type", "application/json; charset=utf-8" )

  dirnamesmap := map[string] string { "sitemap" : "server/data", "factories" : "client/images" }
  vars := mux.Vars(r)
  dirname := "/home/stephen/hacking/fedwiki/" + dirnamesmap[vars["map"]]

  log.Println( dirname )

  slugs, err := ioutil.ReadDir(dirname)

  if err != nil {
    log.Fatal(err)
  }
  
  //items := make(map[string] *MapItem)
  items := [] *MapItem {}

  for _, value := range slugs {
    if value.IsDir() != true {
      item := new( MapItem )
      item.Slug = value.Name()
      item.Date = value.ModTime().Unix()
      item.Title = value.Name()
      item.Synopsis = "synopsis"

      //items[item.Slug] = item
      items = append(items, item)

      //log.Println( items[item.Slug] )
    }
  }

  enc := json.NewEncoder(w)
  enc.Encode(items)
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/{slug:[a-z0-9-]+}.json", JsonHandler)
    r.HandleFunc("/view/{id:[a-z0-9-]+}", ViewHandler)
    r.HandleFunc("/{fn:[A-Za-z0-9-]+.(css|js|png)}", StaticHandler)
    r.HandleFunc("/js/{fn:[A-Za-z0-9-.]+.(css|js)}", JsHandler)
    r.HandleFunc("/js/{sub}/{fn:[A-Za-z0-9_.-]+.(css|js|png)}", JsSubHandler)
    r.HandleFunc("/plugins/{fn:[A-Za-z0-9_.-]+.(coffee|js|json)}", PluginHandler)
    r.HandleFunc("/plugins/{sub}/{fn:[A-Za-z0-9_.-]+.(coffee|js|json)}", PluginSubHandler)

    r.HandleFunc("/system/{map:(sitemap|factories)}.json", MapHandler)

    http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}
