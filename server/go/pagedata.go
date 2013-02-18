package main

import (
)

type Site struct {
  Domain string
  Data Store
  ClientDirectory string
}

//Store is an interface for Data access, gets/puts json data.
type Store interface {
  Location() string
  Exists( path string ) bool
  Get( path string ) * Content
  Put( path string, content * Content )
}

type FileStore struct {
  Directory string
  DefaultDirectory string
}

type Item struct {
  Type string `json:"type,omitempty"`
  Id string `json:"id,omitempty"`
  Text string `json:"text,omitempty"`
  Title string `json:"title,omitempty"`
}

type Entry struct {
  Type string `json:"type,omitempty"`
  Id string `json:"id,omitempty"`
  Date int64 `json:"date,omitempty"`
  After string `json:"after,omitempty"`
  Item *Item `json:"item,omitempty"`
  Site *string `json:"site,omitempty"`
  Order [] *string `json:"order,omitempty"`
}

type Page struct {
  Directory string
  Default_Directory string
  Plugins_Directory string
  Body * Content
}

type Content struct {
  Title string `json:"title,omitempty"`
  Story [] *Item `json:"story,omitempty"`
  Journal [] *Entry `json:"journal,omitempty"`
}

type MapItem struct {
  Slug string `json:"slug"`
  Title string `json:"title"`
  Date int64 `json:"date"`
  Synopsis string `json:"synopsis"`
}

type FactoryInfo struct {
  Name string `json:"name"`
  Title string `json:"title"`
  Category string `json:"category"`
}

type ViewInfo struct {
  Status string
  Slug string
}

