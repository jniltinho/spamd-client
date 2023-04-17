/*
Package main
spamd-client - Golang Spamd SpamAssassin Client
*/
package main

// StatusCode StatusCode
// StatusMsg  string
// Version    string
// Score      float64
// BaseScore  float64
// IsSpam     bool
// Headers    textproto.MIMEHeader
// Msg        *Msg
// Rules      []map[string]string

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	spamc "github.com/baruwa-enterprise/spamd-client/pkg"
	"github.com/baruwa-enterprise/spamd-client/pkg/response"
	flag "github.com/spf13/pflag"
)

var (
	cfg     *Config
	cmdName string
)

// Config holds the configuration
type Config struct {
	Address        string
	Port           int
	UseTLS         bool
	User           string
	UseCompression bool
	RootCA         string
	Filename       string
	Socket         bool
	SockerPath     string
}

func d(r *response.Response) {
	// log.Println("===================================")
	log.Printf("RequestMethod => %v\n", r.RequestMethod)
	log.Printf("StatusCode => %v\n", r.StatusCode)
	log.Printf("StatusMsg => %v\n", r.StatusMsg)
	log.Printf("Version => %v\n", r.Version)
	log.Printf("Score => %v\n", r.Score)
	log.Printf("BaseScore => %v\n", r.BaseScore)
	log.Printf("IsSpam => %v\n", r.IsSpam)
	log.Printf("Headers => %v\n", r.Headers)
	log.Printf("Msg => %v\n", r.Msg)
	log.Printf("Msg.Header => %v\n", r.Msg.Header)
	log.Printf("Msg.Body => %s", r.Msg.Body)
	log.Printf("Msg.Raw => %s", r.Raw)
	log.Printf("Rules => %v\n", r.Rules)
	log.Println("===================================")
}

func init() {
	cfg = &Config{}
	cmdName = path.Base(os.Args[0])
	flag.StringVarP(&cfg.Address, "host", "H", "127.0.0.1", `Specify Spamd host to connect to.`)
	flag.BoolVarP(&cfg.Socket, "socket", "S", true, `Use Linux socket.`)
	flag.StringVarP(&cfg.SockerPath, "socket-path", "s", "/var/run/clamav/spamd-socket", `Specify socket path.`)

	flag.IntVarP(&cfg.Port, "port", "p", 783, `In TCP/IP mode, connect to spamd server listening on given port`)
	flag.StringVarP(&cfg.User, "user", "u", "clamav", `User for spamd to process this message under.`)
	flag.BoolVarP(&cfg.UseCompression, "use-compression", "z", false, `Compress mail message sent to spamd.`)
	flag.StringVarP(&cfg.Filename, "file-name", "F", "", `Specify filename .eml file to scan.`)
}

func parseAddr(a string, p int) (n string, h string) {
	if strings.HasPrefix(a, "/") {
		n = "unix"
		h = a
	} else {
		n = "tcp"
		if strings.Contains(a, ":") {
			h = fmt.Sprintf("[%s]:%d", a, p)
		} else {
			h = fmt.Sprintf("%s:%d", a, p)
		}
	}
	return
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", cmdName)
	fmt.Fprint(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.ErrHelp = errors.New("")
	flag.CommandLine.SortFlags = false
	flag.Parse()
	ctx := context.Background()
	network, address := parseAddr(cfg.Address, cfg.Port)
	m := []byte("Subject: Hello\r\n\r\nHey there!\r\n")

	if cfg.Filename != "" {
		var err error
		m, err = os.ReadFile(cfg.Filename)
		if err != nil {
			log.Println(err)
			return
		}
	}

	if cfg.Socket {
		network = "unix"
		address = cfg.SockerPath
	}
	c, err := spamc.NewClient(network, address, cfg.User, cfg.UseCompression)
	if err != nil {
		log.Println(err)
		return
	}
	c.EnableRawBody()
	ir := bytes.NewReader(m)
	r, e := c.Symbols(ctx, ir)
	if e != nil {
		log.Println(e)
		return
	}
	d(r)
}
