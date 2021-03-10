package main

import(
    "fmt"
    "utils"
    "cache"
    "manifest"
)

func main() {

    if utils.Exists("go.mod"){
        fmt.Println("Yes")
    }else{
        fmt.Println("No")
    }
    cache.Greet()
    hash, err := manifest.GetHash("main.go")
    if err != nil {
        fmt.Println(err)
    }else{
        fmt.Println(hash)
    }
}
