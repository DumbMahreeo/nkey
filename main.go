package main

/*
#cgo LDFLAGS: -lncurses

#include <ncurses.h>
*/
import "C"

import (
    "fmt"
    "time"
    "os"
)

func addStr(text string, vars ...interface{}) {
    C.addstr(
        C.CString(
            fmt.Sprintf(text, vars...),
        ),
    )
    C.refresh()
}

func capture(doKill bool, toStdout bool) string {

    var output string

    C.initscr()

    C.noecho()
    C.raw()

    occupied := false

    C.refresh()

    Loop:
    for {

        char := C.getch()
        go func() {
            if !occupied {
                occupied = true
                time.Sleep(time.Millisecond)
                C.clear()

                if toStdout {
                    output += "\n"
                }

                occupied = false
            }
        }()

        var stringChar string

        switch char {
        case 3:
            addStr("%d\tKill\n\r", char)
            if doKill {
                break Loop
            }
            continue Loop

        case 10:
            stringChar = "Enter"

        case 27:
            stringChar = "Control char"

        case 32:
            stringChar = "Space"

        case 127:
            stringChar = "Backspace"

        default:
            stringChar = string(char)
        }

        addStr("%d\t%s\n\r", char, stringChar)

        if toStdout {
            output += fmt.Sprintf("%d\t%s\n", char, stringChar)
        }
    }

    C.endwin()

    return output

}

func main() {
    doKill := true
    toStdout := false

    for _, arg := range os.Args {
        switch arg {
        case "--help", "-h":
            fmt.Println(
                "Help message\n\n"+
                "--help -h\tShow this help message\n"+
                "--no-kill -k\tPrevent ctrl+c from closing the application\n"+
                "--stdout -s\tPrints the data to stdout (will override --no-kill)\n"+
                "\n\nNote: ctrl+c is equivalent to 3",
            )
            return

        case "--stdout", "-s":
            toStdout = true

        case "--no-kill", "-k":
            doKill = false
        }

    }

    if toStdout {
        doKill = true
    }

    output := capture(doKill, toStdout)

    if toStdout {
        fmt.Printf("Kill on ctrl+c: %t\nPrint to stdout: %t\n\n", doKill, toStdout)
        fmt.Print(output)
    }
}
