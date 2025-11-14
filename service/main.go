package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

	v1 "github.com/suitcase/butler/api/v1"
	"github.com/suitcase/butler/db"
)

func init() {
	fmt.Printf("       ____  _   _ _____ _     _____ ____        \n")
	fmt.Printf("      | __ )| | | |_   _| |   | ____|  _ \\       \n")
	fmt.Printf("      |  _ \\| | | | | | | |   |  _| | |_) |      \n")
	fmt.Printf("      | |_) | |_| | | | | |___| |___|  _ <       \n")
	fmt.Printf("      |____/ \\___/  |_| |_____|_____|_| \\_\\      \n")
	fmt.Printf("                                                 \n")
	fmt.Printf("       github.com/dubcook29/suitcase-butler      \n")
	fmt.Printf("                                                 \n")
	fmt.Printf("           ------  0.0.1(Alpha)  ------          \n")
	fmt.Printf("                                                 \n")
	fmt.Printf("                                                 \n")
}

func main() {
	var (
		dbs     string
		db_host string
		db_port string
		db_user string
		db_pass string
		host    string
		port    int
	)

	flag.StringVar(&dbs, "dbs", "localhost:27017", "connection address for mongodb")
	flag.StringVar(&db_user, "user", "", "username for mongodb connection")
	flag.StringVar(&db_pass, "pass", "", "password formongodb connection")
	flag.StringVar(&host, "h", "localhost", "host / ipaddress for butler api")
	flag.IntVar(&port, "p", 8080, "port number for butler api")
	flag.Parse()

	if dbs_list := strings.Split(dbs, ":"); len(dbs_list) == 2 {
		db_host = dbs_list[0]
		db_port = dbs_list[1]
	} else {
		panic(fmt.Errorf("dbs input error"))
	}

	if err := db.InitialRuntimeDBConnect(context.TODO(), db_host, db_port, db_user, db_pass); err != nil {
		panic(err)
	}

	v1.ButlerAPIServiceStarter(host, strconv.Itoa(port))
}
