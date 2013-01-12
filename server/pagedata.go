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

