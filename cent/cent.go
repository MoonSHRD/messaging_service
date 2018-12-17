package cent

import (
    "fmt"
    "github.com/centrifugal/gocent"
    "log"
    "messaging_service/models"
)

type CentInst struct {
    Client *gocent.Client
}

var centInst = CentInst{nil}

func GetCent() *gocent.Client {
    return centInst.Client
}

func NewCent(c *models.Config) *gocent.Client {
    cConfig :=gocent.Config{
        Addr:fmt.Sprintf("http://%s:%s", c.CHost, c.CPort),
        Key:c.CConfig.Key,
    }
    log.Println(cConfig)
    centInst.Client=gocent.New(cConfig)
    return centInst.Client
}