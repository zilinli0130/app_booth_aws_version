package main

import (
    "fmt"
    "log"
    "net/http" 
    "appstore/handler"  
    "appstore/backend" 
)

func main() {
    fmt.Println("started-service")
    backend.InitElasticsearchBackend()
    backend.InitS3Backend()
    log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}