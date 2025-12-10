# Go bindings for libFTDI

Go bindings for libFTDI library (http://www.intra2net.com/en/developer/libftdi/), providing a Go interface to communicate with FTDI USB devices.

## Overview

This library provides Go bindings for the popular libFTDI library, enabling communication with FTDI USB-to-serial converter chips. It supports various FTDI operations including:

- Device detection and enumeration
- Serial communication
- Bit-bang mode operations
- SPI communication
- EEPROM operations
- LCD control

## Changes Made

This fork includes the following modifications from the original [ziutek/ftdi](https://github.com/ziutek/ftdi) repository, with additional contributions from [mrmikhailv/ftdi](https://github.com/mrmikhailv/ftdi):


- **Updated module path**: Changed from `github.com/ziutek/ftdi` to `github.com/esoutham-lvt/ftdi`
- **Go 1.25 compatibility**: Updated to work with Go 1.25
- **Deprecation handling**: Added wrapper functions for libftdi 1.5 to handle deprecated purge API functions while maintaining backward compatibility with libftdi 1.4
- **Cross-platform build improvements**: Enhanced CGO directives for better Windows, Linux, and macOS support
- **Updated examples**: All example code updated to use the new module path

## Minimum Requirements

### Go Version
- Go 1.25 or later

### System Dependencies

#### Linux
- libftdi1 development package
- pkg-config

```bash
# Ubuntu/Debian
sudo apt-get install libftdi1-dev pkg-config

# CentOS/RHEL/Fedora  
sudo yum install libftdi-devel pkgconfig
# or for newer versions:
sudo dnf install libftdi-devel pkgconfig
```

#### macOS
- libftdi1 (via Homebrew or MacPorts)
- pkg-config

```bash
# Using Homebrew
brew install libftdi pkg-config
```

#### Windows
- Pre-compiled libftdi1 and libusb libraries (included in `libftdi1-1.5/` directory)
- CGO-compatible C compiler (MinGW-w64 recommended)

### Hardware Requirements
- Compatible FTDI USB device (FT232, FT2232, FT4232, etc.)
- Appropriate USB drivers for your FTDI device

## Installation

```bash
go get github.com/esoutham-lvt/ftdi
```

## Usage Examples

The library includes numerous examples in the `examples/` directory:

- `ftdi_chipid/` - Read FTDI chip ID
- `ftdi_dump/` - Dump the eeprom information
- `ftdi_bitbang/` - Bit-bang mode operations
- `ftdi_spi/` - SPI communication
- `ftdi_eeprom/` - EEPROM operations
- And more...

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/esoutham-lvt/ftdi"
)

func main() {
    // Find all FTDI devices
    devices, err := ftdi.FindAll(0, 0)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d FTDI devices\n", len(devices))
    
    // Open first device
    if len(devices) > 0 {
        dev := devices[0]
        err := dev.Open()
        if err != nil {
            log.Fatal(err)
        }
        defer dev.Close()
        
        // Your FTDI operations here...
    }
}
```

## Documentation

For detailed API documentation, visit: [original documentation](https://pkg.go.dev/github.com/ziutek/ftdi?utm_source=godoc)

## License

This library maintains the same license as the original libFTDI library. See the `libftdi1-1.5/Copyright/` directory for complete license information.
