package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/Apakhov/cube/cubeapi/oauth2"
)

var fs = flag.NewFlagSet("cube", flag.ContinueOnError)

var host = fs.String("host", "", "tcp/ip server host, non-empty string")
var port = fs.Int("port", 0, "tcp/ip server port, positive integer")
var token = fs.String("token", "", "your token, non-empty string")
var scope = fs.String("scope", "", "scope of the, token non-empty string")
var secondsToOperate = fs.Int64("sec", 10, "time before request deadline")

func init() {
	fs.StringVar(host, "h", "", "tcp/ip server host, non-empty string")
	fs.IntVar(port, "p", 0, "tcp/ip server port, positive integer")
	fs.StringVar(token, "t", "", "your token, non-empty string")
	fs.StringVar(scope, "s", "", "scope of the token, non-empty string")

	fs.Usage = func() {
		fmt.Println(`Usage of cube:
	cube host port token scope
or with flags:`)

		fs.PrintDefaults()
	}
}

func checkStringFlag(flag *string, flagName string, pos int) int {
	if *flag == "" {
		fs.Set(flagName, fs.Arg(pos))
		if *flag == "" {
			fmt.Println("expected " + flagName)
			os.Exit(-1)
		}
		return pos + 1
	}
	return pos
}

func checkIntFlag(flag *int, flagName string, pos int) int {
	if *flag == 0 {
		fs.Set(flagName, fs.Arg(pos))
		if *flag == 0 {
			fmt.Println("expected " + flagName)
			os.Exit(-1)
		}
		return pos + 1
	}
	return pos
}

const buffLen = 16

func main() {
	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(-1)
	}

	curParam := 0

	curParam = checkStringFlag(host, "host", curParam)
	curParam = checkIntFlag(port, "port", curParam)
	curParam = checkStringFlag(token, "token", curParam)
	checkStringFlag(scope, "scope", curParam)

	address := fmt.Sprintf("%s:%d", *host, *port)

	var err error
	deadline := time.Now().Add(time.Second * time.Duration(*secondsToOperate))

	fmt.Println("connecting to", address)
	conn, err := net.DialTimeout("tcp", address, time.Second*time.Duration(*secondsToOperate))
	if err != nil {
		fmt.Println("failed to dial tcp", err.Error())
		os.Exit(-1)
	}
	defer conn.Close()
	fmt.Println("connected")

	err = conn.SetWriteDeadline(deadline)
	if err != nil {
		fmt.Println("failed to set write deadline", err.Error())
		os.Exit(-1)
	}
	err = conn.SetReadDeadline(deadline)
	if err != nil {
		fmt.Println("failed to set read deadline", err.Error())
		os.Exit(-1)
	}

	fmt.Println("writing")
	buf, err := oauth2.CreateOAUTH2Request(*token, *scope)
	if err != nil {
		fmt.Println("failed to create request", err.Error())
		os.Exit(-1)
	}
	fmt.Println(buf.Bytes(), "buf")
	fmt.Println(string(buf.Bytes()), "str")

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println("failed to write to connection", err.Error())
		os.Exit(-1)
	}

	respBuf := make([]byte, buffLen)
	r := new(oauth2.ResponseOAUTH2)
	response := oauth2.CreateRespBuffer([]byte{})
	go func() {
		response.ParseOAUTH2Resp(r)
		response.End()
	}()
	var readed int
	err = nil

CONNECT_LOOP:
	for {
		select {
		case <-response.WaitChan():
			break CONNECT_LOOP
		default:
			readed, err = conn.Read(respBuf)
			fmt.Printf("pieceLen: %d `%s` %v\n", readed, string(respBuf[:readed]), respBuf[:readed])

			response.Write(respBuf[:readed])
			if err == io.EOF {
				break CONNECT_LOOP
			}
			if err != nil {
				fmt.Println("failed to read from connection", err.Error())
				os.Exit(-1)
			}
		}
	}
	response.Finished()
	<-response.WaitChan()

	err = response.Error()
	if err != nil {
		fmt.Println("failed to parse response", err.Error())
		os.Exit(0)
	}
	fmt.Println(r.String())
	fmt.Printf("%+v", *r)
}
