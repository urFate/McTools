package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Lukaesebrot/mojango"
)

func User() {
	client := mojango.New()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	uuid, err := client.FetchUUID(flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}

	profile, err := client.FetchProfile(uuid, true)
	if err != nil {
		log.Fatalln(err)
	}

	hist, err := client.FetchNameHistory(uuid)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("┏━━━━━━━━━━━━━━━ Profile ━━━━━━━━━━━━━━━\n┃ Name: %v\n┃ UUID: %v\n┣━━━━━━━━━━━━━ Name History ━━━━━━━━━━━━",
		profile.Name, profile.UUID)

	for i, v := range hist {
		var changeTime = ""
		if v.ChangedToAt != 0 {
			changeTime = time.Unix(0, v.ChangedToAt*int64(time.Millisecond)).Format("(2 January 2006 15:04)")
		}

		fmt.Printf("\n┃ %v. %v %v", i+1, v.Name, changeTime)
	}
	fmt.Println("\n┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
