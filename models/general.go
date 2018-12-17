package models

type CentrifugoConfig struct {
    Key string `json:"api_key"`
}

type Config struct {
    Host string
    Port string
    // centrifugo config & host & port
    CConfig CentrifugoConfig
    CHost string
    CPort string
}

type Message struct {
    From string `json:"from"`
    To string `json:"to"`
    Text string `json:"text"`
}