/*
CONTACTBOOK TABLE

Topics:
- An anonymous struct slice.
- Tabular data formatting.
- XML and JSON marshaling.

@author davorpatech
@since 2021-05-26
*/

package main

import "encoding/json"
import "encoding/xml"
import "fmt"

import "strings"

const ttyw = 41

func main() {
  
  contacts := []struct {
      Name string  `json:"name"  xml:"Name"`
      City string  `json:"city,omitempty"   xml:"City,omitempty"`
      Age  int     `json:"age"  xml:"age,attr"`
//      XMLName xml.Name  `xml:"Contact"`
    } /* `xml:"Contacts>Contact"` */ {
      {"James",         "Austin",      42},
      {"Peter Jr.",     "Miami",       33},
      {"Sarah",         "Paris",       21},
      {"Caleb",         "London",      18},
      {"Cindy M.",      "Tokio",       67},
      {"Jan Markus",    "Stuttgart",    9},
      {"Robin Wels",    "Oslo",        39},
      {"AA Mom",        "Madrid",      77},
      {"AA Daddy",      "Madrid",      78},
      {"Jim Padwin",    "Berlin",      18},
      {"Hanna Robbins", "Sommerville", 16},
      {"Antón Pileiro", "Ourense",     37},
    }
  size := len(contacts)
  
  
  // Tabular print
  fmt.Println(Center("MY CONTACT BOOK", ttyw-1))
  fmt.Println(strings.Repeat("=", ttyw))
  // negative values aligns to left, positives to right
  w1,w2,w3,w4 := len(fmt.Sprint(size)),-13,3,-14
  w4 += w1
  fmt.Printf(
          `%*v#  %*v | %*v | %*v |
`,
          w1,"#", w2,"NAME", w3,"AGE", w4,"CITY")
  fmt.Println(strings.Repeat("=", ttyw))
  switch size {
    case 0:
      fmt.Println()
      fmt.Println(Center("NO DATA FOUND", ttyw))
      fmt.Println()
    default: for i,contact := range contacts {
      fmt.Printf(
          `%0*d#  %*v | %*d | %*v |
`,
          w1, i+1,
          w2, contact.Name,
          w3, contact.Age,
          w4, contact.City)
    }
  }
  fmt.Println(strings.Repeat("=", ttyw))
  
  //
  // XML ENCODING
  // - Attributes to include must be exportable.
  // - Needs a parent structure defining root element
  //
  root := struct {
      XMLName xml.Name  `xml:"Root"`
      Data interface{}  `xml:"Contacts>Contact"`
    }{
      Data: contacts,
    }
  fmt.Println()
  fmt.Println(strings.Repeat("=", ttyw))
  fmt.Printf(Center("XML", ttyw))
  fmt.Println(strings.Repeat("=", ttyw))
  fmt.Println()
  if b, err := xml.MarshalIndent(root, "", "  "); err != nil {
    fmt.Printf("error: %v\n", err)
  } else {
    fmt.Print(xml.Header, string(b), "\n")
  }
  
  
  //
  // JSON ENCODING
  // - Attributes to include must be exportable
  //
  fmt.Println()
  fmt.Println(strings.Repeat("=", ttyw))
  fmt.Printf(Center("JSON", ttyw))
  fmt.Println(strings.Repeat("=", ttyw))
  fmt.Println()
  if b, err := json.MarshalIndent(contacts, "", "  "); err != nil {
    fmt.Printf("error: %v\n", err)
  } else {
    fmt.Println(string(b))
  }
}


// Center tries to center a string to a given width using whitespace as filling characters.
// It preserves codepoints joined when distributable width is computed.
func Center(s string, w int) (cs string) {
    if len(s) == 0 {
        cs = strings.Repeat(" ", w)
        return 
    }
    scp := []rune(s)
//    sc = fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w + len(s))/2, s))
    cs = fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w + len(scp))/2, s))
    return
}
