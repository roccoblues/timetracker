package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const dbFile = "tt.json"
const roundTo = 15 * time.Minute

func main() {
	db := db{}

	if _, err := os.Stat(dbFile); err == nil {
		data, err := ioutil.ReadFile(dbFile)
		if err != nil {
			fmt.Printf("Failed to read file '%s': %v\n", dbFile, err)
			os.Exit(1)
		}

		err = db.decode(data)
		if err != nil {
			fmt.Printf("Failed to decode: %v\n", err)
			os.Exit(1)
		}
	}

	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "start":
			err := db.addStartTime(time.Now())
			if err != nil {
				fmt.Printf("Failed to start timer: %v\n", err)
				os.Exit(1)
			}
		case "stop":
			err := db.addStopTime(time.Now())
			if err != nil {
				fmt.Printf("Failed to stop timer: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Printf("Unknown command: %s\n", cmd)
			os.Exit(1)
		}
	}

	if db.Modified {
		data, err := db.encode()
		if err != nil {
			fmt.Printf("Failed to encode: %v\n", err)
			os.Exit(1)
		}

		err = ioutil.WriteFile(dbFile, data, 0644)
		if err != nil {
			fmt.Printf("Failed to write file '%s': %v\n", dbFile, err)
			os.Exit(1)
		}
	}

	db.print(os.Stdout)
}
