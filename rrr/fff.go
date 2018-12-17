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
    "github.com/libp2p/go-libp2p-peerstore"
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
    dest := flag.String("d", "", "Destination multiaddr string")
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
    
    fmt.Println("This node's multiaddresses:")
    for _, la := range host.Addrs() {
        fmt.Printf(" - %v\n", la)
    }
    fmt.Println()
    
    // Turn the destination into a multiaddr.
    maddr, err := multiaddr.NewMultiaddr(*dest)
    if err != nil {
        log.Fatalln(err)
    }
    
    // Extract the peer ID from the multiaddr.
    info, err := peerstore.InfoFromP2pAddr(maddr)
    if err != nil {
        log.Fatalln(err)
    }
    log.Println(info)
    
    // Add the destination's peer multiaddress in the peerstore.
    // This will be used during connection and stream creation by libp2p.
    host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
    
    // Start a stream with the destination.
    // Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
    s, err := host.NewStream(context.Background(), info.ID, "/chat/1.0.0")
    if err != nil {
        panic(err)
    }
    
    //dht.
    idht,err:=dht.New(ctx,host)
    if err != nil {
        log.Println(err)
    }
    idht.Bootstrap(ctx)
    idht.PutValue(ctx,"value",[]byte("value"))
    val,err:=idht.GetValue(ctx,"value")
    if err != nil {
        log.Println(err)
    }
    log.Println(string(val))
    
    // Create a buffered stream so that read and writes are non blocking.
    rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
    
    // Create a thread to read and write data.
    go writeData(rw)
    go readData(rw)
    
    // Hang forever.
    // Hang forever
    <-make(chan struct{})
    
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
