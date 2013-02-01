package main

import (
)

func (c * Content) Synopsis () string {
  text := ""

  if len( c.Story ) > 3 {
    for _, para := range c.Story[:2] {
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

func (c * Content) AddItem ( e * Entry ) {
  if e.After != "" {
    var story [] * Item

    for _, x := range c.Story {
      story = append( story, x )
      if x.Id == e.After {
        story = append( story, e.Item )
      }
    }

    c.Story = story
  } else {
    c.Story = append( c.Story, e.Item )
  }
  c.Journal = append( c.Journal, e )
}

func (c * Content) AddEntry ( e * Entry ) {
  c.Journal = append( c.Journal, e )
}

func (c * Content) ReplaceItem ( e * Entry ) {  
  var story [] * Item

  for _, x := range c.Story {
    if x.Id == e.Id {
      story = append( story, e.Item )
    } else {
      story = append( story, x )
    }
  }

  c.Story = story
  c.Journal = append( c.Journal, e )
}

func (c * Content) Reorder ( e * Entry ) {
  paragraphs := make( map[string] * Item )

  for _, x := range c.Story {
    paragraphs[x.Id] = x 
  }
  
  story := [] * Item {}

  for _, x := range e.Order {
    story = append( story, paragraphs[*x] )
  }

  c.Story = story
  c.Journal = append( c.Journal, e )
}

