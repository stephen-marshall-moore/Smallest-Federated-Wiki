package main

import (
  "encoding/json"
  "log"
  "os"
  "path"
  "strconv"
)

func (p * Page) assert () {
  fi, err := os.Stat( p.Directory )
  if err != nil || !fi.IsDir() {
    log.Fatal( err, "Where is directory, " + p.Directory + "?" )
  }
  fi, err = os.Stat( p.Default_Directory )
  if err != nil || !fi.IsDir() {
    log.Fatal( err, "Where is default directory, " + p.Default_Directory + "?" )
  }
  fi, err = os.Stat( p.Plugins_Directory )
  if err != nil || !fi.IsDir() {
    log.Fatal( err, "Where is directory, " + p.Plugins_Directory + "?" )
  }
}  

func (p * Page) Synopsis () string {
  text := ""

  if len( p.Body.Story ) > 3 {
    for _, para := range p.Body.Story[:2] {
      if para.Type == "paragraph" {
        text += para.Text
      }
    }
  }
  
  if len(text) > 0 {
    return text
  }

  return "A story that does not start with text." 
}

func (p * Page) Read( slug string ) {
  //file, err := os.Open("/home/stephen/hacking/fedwiki/server/data/" + slug) // For read access.
  file, err := os.Open(path.Join( p.Directory,  slug)) // For read access.
  defer file.Close()

  if err != nil {
    file, err = os.Open(path.Join( p.Default_Directory,  slug))
    if err != nil {
      log.Fatal(err)
    } 
  }
  
  enc := json.NewDecoder(file)
  p.Body = new(Content)
  enc.Decode(p.Body)
}

func (p * Page) Create( slug string ) {
  //fn := "/home/stephen/hacking/fedwiki/server/data/" + slug

  file, err := os.Create(path.Join( p.Directory,  slug))
  defer file.Close()

  if err != nil {
    log.Fatal(err)
  }

  enc := json.NewDecoder(file)
  p.Body = new(Content)
  enc.Decode(p.Body)
}

func (p * Page) Write( slug string ) {
  //fn := "/home/stephen/hacking/fedwiki/server/data/" + slug
  fn := path.Join( p.Directory, slug )
  err2 := os.Rename(fn, path.Join( fn , "_" + strconv.Itoa( len( p.Body.Journal ) )))
  if err2 != nil {
    log.Fatal(err2)
  }

  file, err := os.Create(fn)
  defer file.Close()

  if err != nil {
    log.Fatal(err)
  }

  enc := json.NewEncoder(file)
  enc.Encode(p.Body)
}

func (p * Page) AddItem ( e * Entry ) {
  if e.After != "" {
    var story [] * Item

    for _, x := range p.Body.Story {
      story = append( story, x )
      if x.Id == e.After {
        story = append( story, e.Item )
      }
    }

    p.Body.Story = story
  } else {
    p.Body.Story = append( p.Body.Story, e.Item )
  }
  p.Body.Journal = append( p.Body.Journal, e )
}

func (p * Page) AddEntry ( e * Entry ) {
  p.Body.Journal = append( p.Body.Journal, e )
}

func (p * Page) ReplaceItem ( e * Entry ) {  
  var story [] * Item

  for _, x := range p.Body.Story {
    if x.Id == e.Id {
      story = append( story, e.Item )
    } else {
      story = append( story, x )
    }
  }

  p.Body.Story = story
  p.Body.Journal = append( p.Body.Journal, e )
}

func (p * Page) Reorder ( e * Entry ) {
  paragraphs := make( map[string] * Item )

  for _, x := range p.Body.Story {
    paragraphs[x.Id] = x 
  }
  
  story := [] * Item {}

  for _, x := range e.Order {
    story = append( story, paragraphs[*x] )
  }

  p.Body.Story = story
  p.Body.Journal = append( p.Body.Journal, e )
}

