package main

import (
  "encoding/json"
  "log"
  "os"
  "path"
  "strconv"
)

func (fs FileStore) Location() string {
  return fs.Directory
}

func (fs FileStore) Exists( name string ) bool {
  fi, err := os.Stat(path.Join( fs.Directory,  name))

  if err != nil {
    log.Println( "FileStore Exists problem on (" + fs.Directory + ", " + name + "): ", err )
    return false
  }
  
  if fi.IsDir() {
    return false
  }
  return true
}

func (fs FileStore) Get( name string ) * Content {
  file, err := os.Open(path.Join( fs.Directory, "pages",  name))
  defer file.Close()

  if err != nil {
    file, err = os.Open(path.Join( fs.DefaultDirectory,  name))
    if err == nil { 
      enc := json.NewDecoder(file)
      c := new(Content)
      enc.Decode(c)

      fs.Put( name, c )
      return c
    } else {
      log.Println( "FileStore Get problem on (" + fs.DefaultDirectory + ", " + name + "): ", err )
      return nil
    }
  }
  
  enc := json.NewDecoder(file)
  c := new(Content)
  enc.Decode(c)
  return c
}

func (fs FileStore) Put( name string, content * Content ) {
  fn := path.Join( fs.Directory, "pages", name )
  err := os.Rename(fn, fn + "_" + strconv.Itoa( len( content.Journal ) ))
  if err != nil {
      log.Println( "FileStore Put [rename] problem on (" + fs.Directory + "/pages/, " + name + "): ", err )
  }

  file, err := os.Create(fn)
  defer file.Close()

  if err != nil {
      log.Println( "FileStore Put [create] problem on (" + fs.Directory + ", " + name + "): ", err )
  }

  enc := json.NewEncoder(file)
  enc.Encode(content)
}

