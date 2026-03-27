package main

import (
    "organiser/scanner"
    "fmt"
)

func main() {
    dir, err := scanner.Scan()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Stored directory:", dir)

}