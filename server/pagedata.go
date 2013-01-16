package main

import (
)

type Item struct {
  Type string `json:"type"`
  Id string `json:"id"`
  Text string `json:"text"`
}

type Entry struct {
  Type string `json:"type"`
  Id string `json:"id"`
  Date int64 `json:"date"`
  After string `json:"after"`
  Item *Item `json:"item,omitempty"`
  Site *string `json:"site,omitempty"`
}

type Page struct {
  Title string `json:"title"`
  Story [] *Item `json:"story"`
  Journal [] *Entry `json:"journal"`
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

