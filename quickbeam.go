/**

Simple in memory keyed byte store

put<fDelim><key><fDelim><data><mDelim>  -   put data into the store at key
get<fDelim><key><fDelim><mDelim>        -   get data from store at key. supports * for all.
del<fDelim><key><fDelim><mDelim>        -   delete data from store at key. supports * for all.
con<fDelim><key><fDelim><mDelim>        -   consume (get and delete) data from store at key. supports * for all.
cnt<fDelim><mDelim>                     -   show store count in log
det<fDelim><mDelim>                     -   detail entire store in the log
log<fDelim><mDelim>                     -   toggle logging

**/
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net"
)

type message struct {
	action string
	key    []byte
	msg    []byte
	reply  chan []byte
}

func main() {

	version := "0.1.7"
	addr := flag.String("addr", ":12345", "TCP addr to server <host>:<port>")
	inQ := flag.Int("inQ", 1024, "Inbound queue length")
	flag.Parse()

	fDelim := '\t'
	mDelim := '\n'
	log := false

	fmt.Println("quickbeam", version, *addr, mDelim, fDelim, *inQ)

	inbound := make(chan message, *inQ)

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Println("Error serving TCP", err)
		return
	}

	// connections
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go connection(conn, inbound, byte(mDelim), byte(fDelim))
		}
	}()

	data := make(map[string][]byte)

	// read inbound
	for m := range inbound {

		if log {
			fmt.Println(string(m.action), string(m.key), string(m.msg))
		}

		key := string(m.key)

		switch m.action {
		case "put":
			data[key] = m.msg
		case "get":
			if key == "*" {
				for _, v := range data {
					m.reply <- v
				}
			} else {
				msg, _ := data[key]
				m.reply <- msg
			}
		case "del":
			if key == "*" {
				data = make(map[string][]byte)
			} else {
				delete(data, key)
			}
		case "con":
			if key == "*" {
				for _, v := range data {
					m.reply <- v
				}
				data = make(map[string][]byte)
			} else {
				msg, _ := data[key]
				m.reply <- msg
				delete(data, key)
			}
		case "cnt":
			fmt.Println(len(data))
		case "det":
			for k, v := range data {
				fmt.Printf("%v%v%v%v", k, string(fDelim), string(v), string(mDelim))
			}
		case "log":
			log = !log
			if log {
				fmt.Println("Logging turned on")
			}
		}

	}

}

// connection reads bytes off the client conn, parses the action and feeds the bytes into main
// along with returning any responses
func connection(conn net.Conn, in chan message, mDelim, fDelim byte) {

	var d []byte
	d = append(d, fDelim)
	r := make(chan []byte, 1024)
	reader := bufio.NewReader(conn)

	// write the responses back to the connection
	go func() {
		for m := range r {
			m = append(m, mDelim)
			_, err := conn.Write(m)
			if err != nil {
				break
			}
		}
	}()

	// read
	for {

		line, err := reader.ReadBytes(mDelim)
		if err != nil {
			break
		}

		// parse bytes
		fields := bytes.Split(line, d)

		var msg []byte

		if len(fields) == 0 {
			continue
		} else if len(fields) >= 3 {
			msg = fields[2]
			if len(fields) >= 4 {
				// cater for delimiters in the message body by appending the fields
				for _, f := range fields[3:] {
					msg = append(msg, fDelim)
					msg = append(msg, f...)
				}
			}
		}

		action := string(fields[0])
		if action != "put" &&
			action != "get" &&
			action != "del" &&
			action != "con" &&
			action != "det" &&
			action != "log" &&
			action != "cnt" {
			continue
		}

		in <- message{action: action, key: fields[1], msg: msg, reply: r}
	}
}
