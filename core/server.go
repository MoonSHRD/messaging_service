package core

import (
    "context"
    "encoding/json"
    "log"
    "messaging_service/cent"
    "messaging_service/models"
    "net/http"
)

type Server struct {
    Config *models.Config
}

func (s Server) Serve() {
    ctx := context.Background()
    centClient:=cent.GetCent()
    channels,err:=centClient.Channels(ctx)
    if err != nil {
        log.Fatalln(err)
    }
    log.Println("channels:",channels)
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        var msg models.Message
        err := json.NewDecoder(r.Body).Decode(&msg)
        if err != nil {
            log.Println(err)
            return
        }
        log.Println(msg)
        //log.Println("got sth")
        //text:=r.Body
        //log.Println(text)
        //log.Println(r.ParseForm())
        //log.Println(r)
        //w.Write([]byte("Ok"))
    })
    
    log.Println("Server starting at " + s.Config.Host+":"+s.Config.Port)
    err=http.ListenAndServe(s.Config.Host+":"+s.Config.Port, nil)
    if err != nil {
        log.Fatalln(err)
    }
}
