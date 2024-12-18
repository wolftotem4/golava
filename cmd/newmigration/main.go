package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const DIR = "./database/migrations"

func main() {
	flag.Parse()

	name := flag.Arg(0)

	if name == "" {
		fmt.Println("Please provide a name for the migration")
		return
	}

	timestamp := time.Now().Unix()
	upFilename := fmt.Sprintf("%d_%s.up.sql", timestamp, name)
	downFilename := fmt.Sprintf("%d_%s.down.sql", timestamp, name)

	fmt.Printf("Creating migration files: %s, %s\n", upFilename, downFilename)

	upFile, err := os.Create(fmt.Sprintf("%s%c%s", DIR, os.PathSeparator, upFilename))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer upFile.Close()

	downFile, err := os.Create(fmt.Sprintf("%s%c%s", DIR, os.PathSeparator, downFilename))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer downFile.Close()

	fmt.Fprintln(upFile, "-- +migrate Up")
	fmt.Fprintln(upFile, "-- SQL in section 'Up' is executed when this migration is applied")
	fmt.Fprintln(upFile, "")
	fmt.Fprintln(downFile, "-- +migrate Down")
	fmt.Fprintln(downFile, "-- SQL section 'Down' is executed when this migration is rolled back")
	fmt.Fprintln(downFile, "")

	fmt.Println("Migration files created successfully")
}
