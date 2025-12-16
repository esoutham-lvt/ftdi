package main

import (
	"flag"
	"fmt"
	"os"
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
	vendor  = flag.Int("vendor", 0x0403, "PCI vendor id")
	product = flag.Int("product", 0x6001, "PCI product id")
	fix     = flag.Bool("fix", false, "fix CBUS2 and CBUS3 as IO")
	filter  = flag.String("filter", "FT232R", "Filter chips which must contain this description")
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

		e := d.EEPROM()

		// if fix option is set, modify the EEPROM to set CBUS2 and CBUS3 as IO
		if *fix == true {
			checkErr(e.Read())
			checkErr(e.Decode())
			// Set CBUS2 and CBUS3 to "CBUS I/O"
			e.SetCBusFunction(2, ftdi.CBusIOMode)
			e.SetCBusFunction(3, ftdi.CBusIOMode)
			checkErr(e.Build())
			checkErr(e.Write())
			fmt.Println("  Fixed CBUS2 and CBUS3 to CBUS I/O")
		}

		// Re-Read
		checkErr(e.Read())
		checkErr(e.Decode())

		// Display Selective EEPROM contents
		fmt.Printf("  Manufacturer: %s\n", e.ManufacturerString())
		fmt.Printf("  Product: %s\n", e.ProductString())
		fmt.Printf("  Serial: %s\n", e.SerialString())
		fmt.Printf("    CBus[2]: %s\n", e.CBusFunction(2).String())
		fmt.Printf("    CBus[3]: %s\n", e.CBusFunction(3).String())
	}

}
