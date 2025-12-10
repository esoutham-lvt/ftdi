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
	vendor   = flag.Int("vendor", 0x0403, "PCI vendor id")
	product  = flag.Int("product", 0x6001, "PCI product id")
	dump     = flag.Bool("dump", false, "Dump specific EEPROM contents")
	fullDump = flag.Bool("fulldump", false, "Dump full EEPROM contents")
	filter   = flag.String("filter", "", "Filter devices by description substring")
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

		if *dump == true || *fullDump == true {
			d, err := ftdi.OpenUSBDev(devInfo, ftdi.ChannelAny)
			checkErr(err)
			defer d.Close()

			e := d.EEPROM()
			checkErr(e.Read())
			checkErr(e.Decode())
			if *fullDump == true {
				fmt.Println(e) //Full dump
			} else {
				//Selective dump
				fmt.Printf("  Manufacturer: %s\n", e.ManufacturerString())
				fmt.Printf("  Product: %s\n", e.ProductString())
				fmt.Printf("  Serial: %s\n", e.SerialString())
				//fmt.Printf("  VendorID: 0x%04X\n", e.VendorId())
				//fmt.Printf("  ProductID: 0x%04X\n", e.ProductId())
				//fmt.Printf("  SelfPowered: %t\n", e.SelfPowered())
				fmt.Printf("    CBus[2]: %s\n", e.CBusFunction(2).String())
				fmt.Printf("    CBus[3]: %s\n", e.CBusFunction(3).String())
			}
		}
	}

}
