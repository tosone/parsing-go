package main

import (
	"smartconn.cc/tosone/ra-plus/newparsing/follow"
	"smartconn.cc/tosone/ra-plus/store"

	_ "smartconn.cc/tosone/ra-plus/newconfig"
)

func main() {
	if err := store.Initialize(); err != nil {
		panic(err)
	}
	if err := follow.GetBase(); err != nil {
		panic(err)
	}
}
