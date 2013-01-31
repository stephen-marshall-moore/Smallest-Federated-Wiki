package main

import (
  "encoding/json"
  "log"
  "os"
  "path"
  "strconv"
)

func (fs FileStore) Exists( path string ) bool {
  fi, err := os.Stat(path.Join( fs.Directory,  path))

  if err != nil {
    log.Println( "FileStore Exists problem on (" + fs.Directory + ", " + path + "): " + err )
    return false
  }
  
  if fi.IsDir() {
    return false
  }
  return true
}

func (fs FileStore) Get( path string ) * Content {
  file, err := os.Open(path.Join( fs.Directory,  path))
  defer file.Close()

  if err != nil {
    file, err = os.Open(path.Join( p.Default_Directory,  path))
    if err == nil { 
      enc := json.NewDecoder(file)
      c := new(Content)
      enc.Decode(c)

      fs.Put( path, c )
      return c
    } else {
      log.Println( "FileStore Get problem on (" + fs.Default_Directory + ", " + path + "): " + err )
      return nil
    }
  }
  
  enc := json.NewDecoder(file)
  c := new(Content)
  enc.Decode(c)
  return c
}

func (fs FileStore) Put( path string, content * Content ) {
  fn := path.Join( p.Directory, path )
  err2 := os.Rename(fn, path.Join( fn , "_" + strconv.Itoa( len( content.Journal ) )))
  if err2 != nil {
      log.Println( "FileStore Put [rename] problem on (" + fs.Directory + ", " + path + "): " + err )
  }

  file, err := os.Create(fn)
  defer file.Close()

  if err != nil {
      log.Println( "FileStore Put [create] problem on (" + fs.Directory + ", " + path + "): " + err )
  }

  enc := json.NewEncoder(file)
  enc.Encode(content)
}

