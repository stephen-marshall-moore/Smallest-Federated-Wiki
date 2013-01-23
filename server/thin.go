package main

import (
  //"fmt"
  "encoding/json"
  //"io/ioutil"
  "log"
  "net/http"
  //"net/url"
  "os"
  "path"
  "path/filepath"
  "regexp"
  "strconv"
  "strings"
  "text/template"
  "github.com/gorilla/mux"
)

func Synopsis ( page Page ) string {
  text := ""

  if len( page.Story ) > 3 {
    for _, p := range page.Story[:2] {
      if p.Type == "paragraph" {
        text += p.Text
      }
    }
  }
  
  if len(text) > 0 {
    return text
  }

  return "A story that does not start with text." 
}

func (p * Page) Read( slug string ) {

  file, err := os.Open("/home/stephen/hacking/fedwiki/server/data/" + slug) // For read access.
  defer file.Close()

  if err != nil {
    log.Fatal(err)
  }

  enc := json.NewDecoder(file)
  enc.Decode(p)
}

func (p * Page) Create( slug string ) {
  fn := "/home/stephen/hacking/fedwiki/server/data/" + slug

  file, err := os.Create(fn)
  defer file.Close()

  if err != nil {
    log.Fatal(err)
  }

  enc := json.NewEncoder(file)
  enc.Encode(p)
}

func (p * Page) Write( slug string ) {
  fn := "/home/stephen/hacking/fedwiki/server/data/" + slug
  
  err2 := os.Rename(fn, fn + "_" + strconv.Itoa( len( p.Journal ) ))
  if err2 != nil {
    log.Fatal(err2)
  }

  file, err := os.Create(fn)
  defer file.Close()

  if err != nil {
    log.Fatal(err)
  }

  enc := json.NewEncoder(file)
  enc.Encode(p)
}

func (p * Page) AddItem ( e * Entry ) {
  if e.After != "" {
    var story [] * Item

    for _, x := range p.Story {
      story = append( story, x )
      if x.Id == e.After {
        story = append( story, e.Item )
      }
    }

    p.Story = story
  } else {
    p.Story = append( p.Story, e.Item )
  }
  p.Journal = append( p.Journal, e )
}

func (p * Page) AddEntry ( e * Entry ) {
  p.Journal = append( p.Journal, e )
}

func (p * Page) ReplaceItem ( e * Entry ) {  
  var story [] * Item

  for _, x := range p.Story {
    if x.Id == e.Id {
      story = append( story, e.Item )
    } else {
      story = append( story, x )
    }
  }

  p.Story = story
  p.Journal = append( p.Journal, e )
}

func (p * Page) Reorder ( e * Entry ) {
  paragraphs := make( map[string] * Item )

  

  for _, x := range p.Story {
    paragraphs[x.Id] = x 
  }
  
  story := [] * Item {}

  for _, x := range e.Order {
    story = append( story, paragraphs[*x] )
  }

  p.Story = story
  p.Journal = append( p.Journal, e )
}

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
  defer file.Close()

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



func MultiViewHandler ( w http.ResponseWriter, r *http.Request ) {
  //fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
  vars := mux.Vars(r)

  log.Println( vars )

  limit, converr := strconv.Atoi(vars["count"])
  if converr != nil { panic(converr) }

  data := make([]*ViewInfo, limit)

  for i := 0; i < limit; i++ {
    datum := new(ViewInfo)
    if i == 0 {
      datum.Status = "active"
    }
    datum.Slug = vars["slug_" + strconv.Itoa(i)]
 
    data[i] = datum
  } 
    
  tmpl, err := template.ParseFiles("/home/stephen/hacking/fedwiki/server/templates/layout.html")
  
  if err != nil { panic(err) }
  /***
  data := struct {
    Title string
    Status string
    Slug string
  } { vars["slug_0"], "active", vars["slug_0"] }
  ***/
  stuff := struct {
    Title string
    Slugs [] * ViewInfo
  } { "my title", data }

  err = tmpl.Execute(w, stuff)
  if err != nil { panic(err) }
}

func SiteMapHandler ( w http.ResponseWriter, r *http.Request ) {
  appRoot := "/home/stephen/hacking/fedwiki"

  w.Header().Set( "Content-Type", "application/json; charset=utf-8" )

  //dirnamesmap := map[string] string { "sitemap" : "server/data", "factories" : "client/images" }
  //vars := mux.Vars(r)
  //dirname := "/home/stephen/hacking/fedwiki/" + dirnamesmap[vars["map"]]

  //log.Println( dirname )

  //slugs, err := ioutil.ReadDir(dirname)

  pattern := appRoot + "/server/data/*"
  slugs, err := filepath.Glob(pattern)

  if err != nil {
    log.Fatal(err)
  }
  
  //items := make(map[string] *MapItem)
  items := [] *MapItem {}

  for _, value := range slugs {
    //if value.IsDir() != true {
      item := new( MapItem )

      info := new( Page )

      file, err := os.Open(value)
      if err != nil {
        log.Fatal(err)
      }

      fi, err2 := os.Stat(value)
      if err2 != nil {
        log.Fatal(err2)
      }

      enc := json.NewDecoder(file)
      enc.Decode(info)

      item.Slug = path.Base(value)
      item.Date = fi.ModTime().Unix()
      item.Title = info.Title
      item.Synopsis = Synopsis(*info)

      items = append(items, item)

      //log.Println( items[item.Slug] )
    //}
  }

  enc := json.NewEncoder(w)
  enc.Encode(items)
}

func FactoriesHandler ( w http.ResponseWriter, r *http.Request ) {
  appRoot := "/home/stephen/hacking/fedwiki"

  w.Header().Set( "Content-Type", "application/json; charset=utf-8" )

  //dirnamesmap := map[string] string { "sitemap" : "server/data", "factories" : "client/images" }
  //vars := mux.Vars(r)
  //dirname := "/home/stephen/hacking/fedwiki/" + dirnamesmap[vars["map"]]

  //log.Println( dirname )

  //slugs, err := ioutil.ReadDir(dirname)

  pattern := appRoot + "/client/plugins/*/factory.json"
  slugs, err := filepath.Glob(pattern)

  if err != nil {
    log.Fatal(err)
  }
  
  //items := make(map[string] *MapItem)
  items := [] *FactoryInfo {}

  for _, value := range slugs {
    //if value.IsDir() != true {
      //item := new( MapItem )

      info := new(FactoryInfo)

      file, err := os.Open(value)
      if err != nil {
        log.Fatal(err)
      }

      enc := json.NewDecoder(file)
      enc.Decode(info)

      //items[item.Slug] = item
      items = append(items, info)

      //log.Println( items[item.Slug] )
    //}
  }

  enc := json.NewEncoder(w)
  enc.Encode(items)
}


func ActionHandler ( w http.ResponseWriter, r *http.Request ) {
  vars := mux.Vars( r )

  log.Println( "action" + ": " + vars["slug"] )
  //log.Println( r )

  err := r.ParseForm()

  if err != nil {
    log.Panic(err)
  }

  vals := r.Form

  jsaction := vals["action"]

  log.Println( jsaction )

  var entry Entry

  dec := json.NewDecoder(strings.NewReader( jsaction[0] ))

  err = dec.Decode(&entry)
  if err != nil {
    log.Panic(err)
  }

  log.Println( entry.Type )
  if entry.Item != nil {
    log.Println( "> " + entry.Item.Text )
  }
  
  page := new(Page)

  switch entry.Type {
    case "add": {
      page.Read( vars["slug"] )
      page.AddItem( &entry )
      page.Write( vars["slug"] ) 
    }
    case "edit": {
      page.Read( vars["slug"] )
      page.ReplaceItem( &entry )
      page.Write( vars["slug"] )
    }
    case "create": {
      if entry.Item != nil {
        page.Title = entry.Item.Title
      }
      page.AddEntry( &entry )
      page.Create( vars["slug"] )
    }
    case "move": {
      page.Read( vars["slug"] )
      page.Reorder( &entry )
      page.Write( vars["slug"] )
    }
  }

  log.Println( page )
}

func MatchMultiples(req *http.Request, m *mux.RouteMatch) bool {
  pat := regexp.MustCompile( "(/[a-zA-Z0-9:.-]+/[a-z0-9-]+(_rev[0-9]+)?)+" )

  flag := pat.MatchString( req.URL.String() ) 

  if flag {
    m.Vars = make(map[string]string)
    for i, match := range strings.Split(req.URL.String()[1:], "/") {
      if match == "page" {
        log.Println( "page in match multiples" )
        return false
      }

      if i % 2 == 0 {
        m.Vars["site_" + strconv.Itoa(i/2)] = match
        m.Vars["count"] = strconv.Itoa((i/2) + 1 )
      } else {
        m.Vars["slug_" + strconv.Itoa(i/2)] = match
      }
    }
    return true
  }

  return false
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/{slug:[a-z0-9-]+}.json", JsonHandler)
    //r.HandleFunc("/view/{id:[a-z0-9-]+}", ViewHandler)
    r.HandleFunc("/{fn:[A-Za-z0-9-]+.(css|js|png)}", StaticHandler)
    r.HandleFunc("/js/{fn:[A-Za-z0-9-.]+.(css|js)}", JsHandler)
    r.HandleFunc("/js/{sub}/{fn:[A-Za-z0-9_.-]+.(css|js|png)}", JsSubHandler)
    r.HandleFunc("/plugins/{fn:[A-Za-z0-9_.-]+.(coffee|js|json)}", PluginHandler)
    r.HandleFunc("/plugins/{sub}/{fn:[A-Za-z0-9_.-]+.(coffee|js|json)}", PluginSubHandler)
    //r.HandleFunc("/{multi:((/?[a-zA-Z0-9:.-]+/[a-z0-9-]+(_rev[0-9]+)?)+)}", MultiViewHandler)
    //r.HandleFunc("/{site:[a-zA-Z0-9:.-]+}/{slug:[a-z0-9-]+(_rev[0-9]+)?}", MultiViewHandler)

    r.HandleFunc("/system/sitemap.json", SiteMapHandler)
    r.HandleFunc("/system/factories.json", FactoriesHandler)

    route7 := r.HandleFunc("/page/{slug:[a-z0-9-]+}/action",ActionHandler)
    route7.Methods("PUT", "GET")

    route := r.MatcherFunc(MatchMultiples)
    route.HandlerFunc(MultiViewHandler)

    http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}
