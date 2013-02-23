package main

import (
  //"fmt"
  "encoding/json"
  //"io"
  "io/ioutil"
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
  "github.com/stephen-marshall-moore/openid.go/src/openid"
)

var appRoot = "/home/fedwiki/smallest"
var dataRoot = "/home/fedwiki/smallest/data"
//var appRoot = "/home/stephen/hacking/Smallest-Federated-Wiki"
//var dataRoot = "/home/stephen/hacking/Smallest-Federated-Wiki/data"
var rootDefault = appRoot + "/default-data"

var baseStore = FileStore {
  Directory: dataRoot + "/farm/wiki.example.com",
  DefaultDirectory: rootDefault }

var base = Site {
  Domain: "wiki.example.com",
  Home: "welcome-visitors",
  Data: baseStore,
  ClientDirectory: appRoot + "/client" }

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
    fname = vars["slug"] //path.Join( site.Data.Location(), "pages", vars["slug"] )
    log.Println( "JsonHandler: " + fname )
  } else {
    fname = "missing-page" //path.Join( rootDefault, "oops.html" )
  }

  //http.ServeFile(w,r,fname)
  content := site.Data.Get( fname )
  //io.WriteString( w, content )
  if content != nil {
    enc := json.NewEncoder(w)
    enc.Encode(content)
  } else {
    http.NotFound(w,r)
  }

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

  loginUrl := r.FormValue( "LoginButton" )

  if site != nil {

    if url, err := openid.RedirectUrl(loginUrl,
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

    siteOwner := site.Owner()
    if siteOwner == nil {
      site.SetOwner(&id)
      session.Values["authenticated"] = 42
      session.Values["id"] = id
      session.Save(r,w)
      http.Redirect(w, r, "/view/welcome-visitors", 302)
    } else {
      if *siteOwner == id {
        session.Values["authenticated"] = 42
        session.Values["id"] = id
        session.Save(r,w)
        http.Redirect(w, r, "/view/welcome-visitors", 302)
      }
    }

    http.Redirect(w, r, "/view/welcome-visitors", 403)

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

func IsAuthenticated ( r * http.Request ) bool {
  session, _ := cookieStore.Get( r, "wiki-woko" )

  return session.Values["authenticated"] == 42
}

func (s Site) Owner() * string {
  file, err := os.Open(path.Join( s.Data.Location(), "status",  "openid.identity"))
  //defer file.Close()

  if err != nil {
      log.Println( "Unowned site: " + s.Domain )
      return nil
  }
  
  content, err2 := ioutil.ReadAll(file)
  if err2 != nil {
    log.Println( "Unowned site: " + s.Domain, err2 )
    return nil
  }

  str := string(content)

  return &str
}

func (s Site) SetOwner(openid * string) {
  err := ioutil.WriteFile(path.Join( s.Data.Location(), "status",  "openid.identity"), [] byte(*openid), 0666)
  //defer file.Close()

  if err != nil {
      log.Println( "Unable to set owner owned site: " + s.Domain + ", " + *openid, err )
  }
}

func WelcomeHandler ( w http.ResponseWriter, r *http.Request ) {
  site := RequestedSite( r )

  if site == nil {
    log.Println( "Welcome Handler: unexpectedly site is nil!" )
    http.Error(w, http.StatusText(403), 403)
  }

  log.Println( "Welcome!" )

  data := make([]*ViewInfo, 1)

  datum := new(ViewInfo)

  datum.Status = "active"
  datum.Slug = site.Home
 
  data[0] = datum
    
  tmpl, err := template.ParseFiles(appRoot + "/server/go/templates/layout.html")
  
  if err != nil { panic(err) }

  stuff := struct {
    Login bool
    Title string
    Slugs [] * ViewInfo
  } { IsAuthenticated(r), "my title", data }

  err = tmpl.Execute(w, stuff)
  if err != nil { panic(err) }
}

func FaviconHandler ( w http.ResponseWriter, r *http.Request ) {
  site := RequestedSite( r )

  if site == nil {
    log.Println( "Favicon Handler: unexpectedly site is nil!" )
    http.Error(w, http.StatusText(403), 403)
  }
  
  iconPath := path.Join( site.Data.Location(), "status", "favicon.png" )

  _, err := os.Stat(iconPath)

  if err == nil {
    http.ServeFile(w,r,iconPath)
  }
  http.NotFound(w,r)
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
    
  tmpl, err := template.ParseFiles(appRoot + "/server/go/templates/layout.html")
  
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

      info := new( Content )
      //info.Body = new( Content )

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
      item.Synopsis = info.Synopsis()

      items = append(items, item)

      //log.Println( items[item.Slug] )
    //}
  }

  enc := json.NewEncoder(w)
  enc.Encode(items)
}

func FactoriesHandler ( w http.ResponseWriter, r *http.Request ) {
  //appHome := appRoot + "/fedwiki"

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

    r.HandleFunc("/{slug:(home|index|)}", WelcomeHandler)
    r.HandleFunc("/favicon.png", FaviconHandler)
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
