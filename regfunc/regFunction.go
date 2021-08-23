package regfunction

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

func RegFunction(w http.ResponseWriter, r *http.Request) {
    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s", reqBody)
    _, err = fmt.Fprintf(w, "receive logs.")
    if err != nil {
        return
    }
}