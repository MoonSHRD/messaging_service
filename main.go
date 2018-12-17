package main

import (
    "bufio"
    "context"
    "flag"
    "fmt"
    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p-crypto"
    "github.com/libp2p/go-libp2p-kad-dht"
    "github.com/libp2p/go-libp2p-net"
    "github.com/multiformats/go-multiaddr"
    "io"
    "log"
    mrand "math/rand"
    "os"
)

// -c /home/nikita/local-chat/centrifugo/config.json -p 8020

func main() {
    ctx := context.Background()
    //confPath := flag.String("c", "./config.json", "centrifugo config.json path")
    sourcePort := flag.Int("p", 8020, "ms port")
    //host := flag.String("h", "localhost", "ms host")
    //cHost := flag.String("ch", "localhost", "centrifugo host")
    //cPort := flag.String("cp", "8000", "centrifugo port")
    flag.Parse()
    var r io.Reader
    r = mrand.New(mrand.NewSource(int64(*sourcePort)))
    prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
    if err != nil {
        panic(err)
    }
    fafa,_:=crypto.MarshalPrivateKey(prvKey)
    log.Println(crypto.ConfigEncodeKey(fafa))
    
    sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort))
    
    host, err := libp2p.New(
        ctx,
        libp2p.ListenAddrs(sourceMultiAddr),
        libp2p.Identity(prvKey),
    )
    if err != nil {
       panic(err)
    }
    log.Println("host id:",host.ID().Pretty())
    log.Println("host addr:",host.Addrs()[2])
    
    host.SetStreamHandler("/chat/1.0.0", handleStream)
    
    // Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
    var port string
    for _, la := range host.Network().ListenAddresses() {
        if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
            port = p
            break
        }
    }
    
    if port == "" {
        panic("was not able to find actual local port")
    }
    
    //fmt.Printf("Run './chat -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, host.ID().Pretty())
    //fmt.Println("You can replace 127.0.0.1 with public IP as well.")
    //fmt.Printf("\nWaiting for incoming connection\n\n")
    
    
    idht,err:=dht.New(ctx,host)
    if err != nil {
        log.Println(err)
    }
    idht.Bootstrap(ctx)
    //idht.
    //idht.Validator
    //namespace:=host.ID().Pretty()
    //idht.
    err=idht.Validator.Validate("/fwafwafwa/2f762f68656c6c6f",[]byte("value"))
    if err != nil {
        log.Println("failed putting data:",err)
    }
    //err=idht.PutValue(ctx,"/"+namespace+"/value",[]byte("value"))
    //if err != nil {
    //    log.Println("failed putting data:",err)
    //}
    //val,err:=idht.GetValue(ctx,"/"+namespace+"/value")
    //if err != nil {
    //    log.Println("failed reading data:",err)
    //} else {
    //    log.Println(string(val))
    //}
    //// Hang forever
    //<-make(chan struct{})
    
    //ctx := context.Background()
    //go func() {
    //   h,err:=libp2p.New(ctx)
    //   if err != nil {
    //       log.Println(err)
    //   }
    //   log.Println(h.ID().Pretty())
    //   log.Println(h.Addrs()[1])
    //   idht,err:=dht.New(ctx,h)
    //   if err != nil {
    //       log.Println(err)
    //   }
    //   idht.Bootstrap(ctx)
    //}()
}


func handleStream(s net.Stream) {
    log.Println("Got a new stream!")
    
    // Create a buffer stream for non blocking read and write.
    rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
    
    go readData(rw)
    go writeData(rw)
    
    // stream 's' will stay open until you close it (or the other side closes it).
}
func readData(rw *bufio.ReadWriter) {
    for {
        str, _ := rw.ReadString('\n')
        
        if str == "" {
            return
        }
        if str != "\n" {
            // Green console colour: 	\x1b[32m
            // Reset console colour: 	\x1b[0m
            fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
        }
        
    }
}

func writeData(rw *bufio.ReadWriter) {
    stdReader := bufio.NewReader(os.Stdin)
    
    for {
        fmt.Print("> ")
        sendData, err := stdReader.ReadString('\n')
        
        if err != nil {
            panic(err)
        }
        
        rw.WriteString(fmt.Sprintf("%s\n", sendData))
        rw.Flush()
    }
    
}
