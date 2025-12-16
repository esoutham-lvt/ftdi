package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/esoutham-lvt/ftdi"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	filter = flag.String("filter", "FT232R", "Filter chips which must contain this description")
)

func main() {
	flag.Parse()
	ftDevList, err := ftdi.FindAll(0x0403, 0x6001) //Get a list of all FTDI parts connected
	checkErr(err)

	fmt.Printf("Found %d FTDI devices:\n", len(ftDevList))
	for i, devInfo := range ftDevList {
		if *filter != "" && !strings.Contains(devInfo.Description, *filter) {
			continue //skip entry
		}

		fmt.Printf(" FTDI Device %d: Mfg=`%s` Serial='%s' Description='%s'\n", i, devInfo.Manufacturer, devInfo.Serial, devInfo.Description)

		// Open the device
		d, err := ftdi.OpenUSBDev(devInfo, ftdi.ChannelAny)
		checkErr(err)
		defer d.Close()

		// Spin in a loop, reading the terminal for an input mask,
		// Exit the loop on 'q' input
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Enter commands (2 character Hex number) (type 'q' to quit):")

		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break // EOF or error
			}

			input := strings.TrimSpace(scanner.Text())
			if strings.Contains(input, "q") {
				fmt.Println("Exiting...")
				break
			}

			fmt.Printf("input: %s\n", input)

			// Add your CBUS manipulation code here
			ioInt, err := strconv.ParseUint(input, 16, 8)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid input. Please enter a valid 2 character Hex number.\n")
				continue
			}

			// Set the CBUS mask
			ioMask := byte(ioInt) // convert uint64 to byte
			fmt.Printf("Set CBUS mask to: 0x%02X\n", ioMask)
			err = d.SetBitmode(ioMask, ftdi.Mode(ftdi.CBusIOMode))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error setting CBUS mask: %v\n", err)
			} else {
				fmt.Println("CBUS mask set successfully.")
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		}
	}
}
