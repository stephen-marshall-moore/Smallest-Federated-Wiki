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
  "github.com/gorilla/sessions"
  "github.com/yohcop/openid.go/src/openid"
)

var rootDefault = "/home/stephen/hacking/fedwiki/default-data"

var baseStore = FileStore {
  Directory: "/home/stephen/hacking/fedwiki/data/farm/wiki.nimbostrati.com",
  DefaultDirectory: rootDefault }

var base = Site {
  Domain: "wiki.nimbostrati.com",
  Data: baseStore,
  ClientDirectory: "/home/stephen/hacking/fedwiki/client" }

var sites = map[string] *Site { base.Domain : &base }

var cookieStore = sessions.NewCookieStore([]byte("not-very-secret"))

// For the demo, we use in-memory infinite storage nonce and discovery
// cache. In your app, do not use this as it will eat up memory and never
// free it. Use your own implementation, on a better database system.
// If you have multiple servers for example, you may need to share at least
// the nonceStore between them.
var nonceStore = &openid.SimpleNonceStore{Store: make(map[string][]*openid.Nonce)}
var discoveryCache = &openid.SimpleDiscoveryCache{}

func RequestedSite ( r * http.Request ) * Site {
  domain := strings.ToLower( strings.Split( r.Host, ":" )[0] )
  log.Println( "Requested Site: " + domain )
  return sites[ domain ]
}

func StaticHandler ( w http.ResponseWriter, r *http.Request ) {
  var fname string

  site := RequestedSite( r )
  if site != nil {
    vars := mux.Vars(r)
    fname = path.Join( site.ClientDirectory, vars["fn"] )
    log.Println( "StaticHandler: " + fname )
  } else {
    fname = path.Join( rootDefault, "oops.html" )
  }

  http.ServeFile(w,r,fname)
}

func JsHandler ( w http.ResponseWriter, r *http.Request ) {
  var fname string

  site := RequestedSite( r )
  if site != nil {
    vars := mux.Vars(r)
    fname = path.Join( site.ClientDirectory, "js", vars["fn"] )
    log.Println( "JsHandler: " + fname )
  } else {
    fname = path.Join( rootDefault, "oops.html" )
  }

  http.ServeFile(w,r,fname)
}

func JsSubHandler ( w http.ResponseWriter, r *http.Request ) {
  var fname string

  site := RequestedSite( r )
  if site != nil {
    vars := mux.Vars(r)
    fname = path.Join( site.ClientDirectory, "js", vars["sub"], vars["fn"] )
    log.Println( "JsSubHandler: " + fname )
  } else {
    fname = path.Join( rootDefault, "oops.html" )
  }

  http.ServeFile(w,r,fname)
}


func PluginHandler ( w http.ResponseWriter, r *http.Request ) {
  var fname string

  site := RequestedSite( r )
  if site != nil {
    vars := mux.Vars(r)
    fname = path.Join( site.ClientDirectory, "plugins", vars["fn"] )
    log.Println( "PluginHandler: " + fname )
  } else {
    fname = path.Join( rootDefault, "oops.html" )
  }

  http.ServeFile(w,r,fname)
}

func PluginSubHandler ( w http.ResponseWriter, r *http.Request ) {
  var fname string

  site := RequestedSite( r )
  if site != nil {
    vars := mux.Vars(r)
    fname = path.Join( site.ClientDirectory, "plugins", vars["sub"], vars["fn"] )
    log.Println( "PluginSubHandler: " + fname )
  } else {
    fname = path.Join( rootDefault, "oops.html" )
  }

  http.ServeFile(w,r,fname)
}


func JsonHandler ( w http.ResponseWriter, r *http.Request ) {
  var fname string

  site := RequestedSite( r )
  if site != nil {
    w.Header().Set( "Content-Type", "application/json; charset=utf-8" )
    vars := mux.Vars(r)
    fname = path.Join( site.Data.Location(), "pages", vars["slug"] )
    log.Println( "JsonHandler: " + fname )
  } else {
    fname = path.Join( rootDefault, "oops.html" )
  }

  http.ServeFile(w,r,fname)
}

func LoginHandler ( w http.ResponseWriter, r *http.Request ) {
  log.Println( "LoginHandler: " )
}

func LogoutHandler ( w http.ResponseWriter, r *http.Request ) {
  session, _ := cookieStore.Get( r, "wiki-woko" )

  session.Values["authenticated"] = 17
  id := session.Values["id"]
 
  log.Println( "LogoutHandler: " + id.(string) )
  session.Save(r,w)

  http.Redirect(w,r, "/view/welcome-visitors", 303)
}


func DiscoverHandler(w http.ResponseWriter, r *http.Request) {
  site := RequestedSite( r )

  if site != nil {

    if url, err := openid.RedirectUrl("https://www.google.com/accounts/o8/id",
      "http://" + r.Host + "/openidcallback",
      ""); err == nil {
      http.Redirect(w, r, url, 303)
    } else {
      log.Print(err)
    }
  } else {
    http.ServeFile(w,r,path.Join( rootDefault, "oops.html" ))
  } 
}

func OpenIdCallbackHandler(w http.ResponseWriter, r *http.Request) {
  site := RequestedSite( r )

  if site == nil {
    http.ServeFile(w,r,path.Join( rootDefault, "oops.html" ))
    return
  } 

  fullUrl := "http://" + r.Host + r.URL.String()
  log.Print(fullUrl)
  id, err := openid.Verify(
      fullUrl,
      discoveryCache, nonceStore)
  if err == nil {

    session, _ := cookieStore.Get( r, "wiki-woko" )
    session.Values["authenticated"] = 42
    session.Values["id"] = id
    session.Save(r,w)

    http.Redirect(w, r, "/view/welcome-visitors", 302)

    /*
    stuff := struct {
      Login bool
      Title string
      Slugs [] * ViewInfo
    } { true, "my title", nil }

    if t, err := template.ParseFiles(path.Join( site.ClientDirectory, "templates", "layout.html")); err == nil {
      t.Execute(w, stuff)
    } else {
      log.Println("WTF")
      log.Print(err)
    }
    */
  } else {
    log.Println("WTF2")
    log.Print(err)
  }
}

/*
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
*/

func IsAuthenticated ( r * http.Request ) bool {
  session, _ := cookieStore.Get( r, "wiki-woko" )

  return session.Values["authenticated"] == 42
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
    Login bool
    Title string
    Slugs [] * ViewInfo
  } { IsAuthenticated(r), "my title", data }

  err = tmpl.Execute(w, stuff)
  if err != nil { panic(err) }
}

func SiteMapHandler ( w http.ResponseWriter, r *http.Request ) {
  //appRoot := "/home/stephen/hacking/fedwiki"
  
  site := RequestedSite( r )

  if site == nil {
  }

  w.Header().Set( "Content-Type", "application/json; charset=utf-8" )

  //dirnamesmap := map[string] string { "sitemap" : "server/data", "factories" : "client/images" }
  //vars := mux.Vars(r)
  //dirname := "/home/stephen/hacking/fedwiki/" + dirnamesmap[vars["map"]]

  //log.Println( dirname )

  //slugs, err := ioutil.ReadDir(dirname)

  pattern := path.Join( site.Data.Location() , "pages", "*" )

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
      info.Body = new( Content )

      file, err := os.Open(value)
      if err != nil {
        log.Fatal(err)
      }

      fi, err2 := os.Stat(value)
      if err2 != nil {
        log.Fatal(err2)
      }

      enc := json.NewDecoder(file)
      enc.Decode(info.Body)

      item.Slug = path.Base(value)
      item.Date = fi.ModTime().Unix()
      item.Title = info.Body.Title
      item.Synopsis = info.Synopsis()

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
  if !IsAuthenticated(r) {
    log.Println( "Action Handler: " + http.StatusText(403) )
    http.Error(w, http.StatusText(403), 403)
    return
  }
  
  site := RequestedSite( r )

  if site == nil {
    log.Println( "Action Handler: unexpectedly site is nil!" )
    http.Error(w, http.StatusText(403), 403)
  }

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
  
  content := new(Content)

  switch entry.Type {
    case "add": {
      content = site.Data.Get( vars["slug"] )
      content.AddItem( &entry )
      site.Data.Put( vars["slug"], content ) 
    }
    case "edit": {
      content = site.Data.Get( vars["slug"] )
      content.ReplaceItem( &entry )
      site.Data.Put( vars["slug"], content )
    }
    case "create": {
      if entry.Item != nil {
        content.Title = entry.Item.Title
      }
      content.AddEntry( &entry )
      site.Data.Put( vars["slug"], content )
    }
    case "move": {
      content = site.Data.Get( vars["slug"] )
      content.Reorder( &entry )
      site.Data.Put( vars["slug"], content )
    }
  }

  log.Println( content )
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

    r.HandleFunc("/login", DiscoverHandler)
    r.HandleFunc("/logout", LogoutHandler)
    r.HandleFunc("/openidcallback", OpenIdCallbackHandler)

    r.HandleFunc("/{slug:[a-z0-9-]+}.json", JsonHandler)
    //r.HandleFunc("/view/{id:[a-z0-9-]+}", ViewHandler)
    //r.HandleFunc("/{fn:[A-Za-z0-9-]+.(css|js|png)}", StaticHandler)
    r.HandleFunc("/{fn:(images/)?[A-Za-z0-9-]+.(css|js|png)}", StaticHandler)
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
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
      log.Fatal( err )
    }
}
