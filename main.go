package main

import "fmt"
import "os"


func main() {
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "usage: %s <TOKEN>\n", os.Args[0])
        os.Exit(1)
    }
}
