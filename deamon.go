package gosync

import (
	"bufio"
	"encoding/gob"
	"flag"
	// "fmt"
	"log"
	"net"
)

func DeamonStart() {
	var lsnHost string
	var lsnPort string
	flag.StringVar(&lsnHost, "h", "", "Please tell me the host ip which you want listen on.")
	flag.StringVar(&lsnPort, "p", "8999", "Please tell me the port which you want listen on.")
	flag.Parse()
	svrln, err := net.Listen("tcp", lsnHost+":"+lsnPort)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := svrln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		dhandleConn(conn)
	}
}

func dhandleConn(conn net.Conn) {
	defer conn.Close()
	cnRd := bufio.NewReader(conn)
	cnWt := bufio.NewWriter(conn)
	dec := gob.NewDecoder(cnRd)
	enc := gob.NewEncoder(cnWt)
	var mg Message
	rcvErr := dec.Decode(&mg)
	if rcvErr != nil {
		log.Println(rcvErr)
	}
	// fmt.Println(mg)

	// **deal with the mg**
	switch mg.MgType {
	case "task":
		hdTask(&mg, cnRd, cnWt, dec, enc)
	case "file":
		hdFile(&mg, cnRd, cnWt, dec, enc)
	default:
		hdNoType(&mg, cnRd, cnWt, dec, enc)
	}
}
