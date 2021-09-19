/*
Dark Theme
Switch statements +10 XP

The user can choose the color theme for the browser:
1. Light
2. Dark
3. Nocturne
4. Terminal
5. Indigo

You are given a program that takes the number as input. Complete the program so that it will output to the console the theme according to input number.

Sample Input
2

Sample Output
Dark

Hint:
Don't forget to use break for each case statement.


@author davorpatech
@since  2021-09-19
*/

package main

import (
    "errors"
    "fmt"
    "io"
    "os"
    "strings"
    "sync"
)

func init() {
    AddTheme("Light")
    AddTheme("Dark")
    AddTheme("Nocturne")
    AddTheme("Terminal")
    AddTheme("Indigo")
    AddTheme("GitHub")
    AddTheme("SoloLearn")
    AddTheme("Dracula")
    AddTheme("Sublime")
}

func main() {
    var theme *Theme

    ask_option:
    for {
        var option string
        PaintMenuTheme()
        fmt.Print("Enter theme option: ")
        _, err := fmt.Scanln(&option)
        fmt.Println()

        if errors.Is(err, io.EOF) {
           return
        }
        if theme = FindTheme(option); theme != nil {
            break ask_option
        }
        if len(option) == 0 {
            return
        }

        fmt.Fprintln(os.Stdout, "ERR!", "Invalid option:", option)
        fmt.Println()
    }

    fmt.Printf("The chosen theme is: %v\n", theme)
}

func PaintMenuTheme() {
    mutex.RLock()
    defer mutex.RUnlock()

    var buf strings.Builder
    mlen := len(fmt.Sprint(lastThemeId))
    for _, theme := range themes {
        buf.WriteString(fmt.Sprintf(
            "%*v. %v\n",
            mlen, theme.GetId(),
            theme.GetName()))
    }

    fmt.Println(buf.String())
}

var (
    lastThemeId uint = 0
    themes = []Theme{}
    mutex    sync.RWMutex
)

func AddTheme(name string) {
    mutex.Lock()
    defer mutex.Unlock()

    lastThemeId++
    theme := Theme{lastThemeId, name}
    themes = append(themes, theme)
}

func FindTheme(option string) *Theme {
    mutex.RLock()
    defer mutex.RUnlock()

    for _, theme := range themes {
        if fmt.Sprint(theme.id) == option {
            return &theme
        }
    }
    return nil
}

type Theme struct{
    id    uint
    name  string
}

func (t Theme) GetId() uint {
    return t.id
}

func (t Theme) GetName() string {
    return t.name
}

/* String implements fmt.Stringer */
func (t Theme) String() string {
    return t.name
}
