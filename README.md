jlog: JSON log tail utility
============================

[![Build Status](https://travis-ci.org/pyr/jlog.svg?branch=master)](https://travis-ci.org/pyr/jlog)

**jlog** allows parsing and easier consumption of JSON serialized
log entries.

## Configuration

**jlog** understands the following command line arguments:

    -n
        Number of lines to tail from a file
    -f
        Watch file for additional output
    -B
	    Force color output (defaults to TTY based behavior)
	-b
	    Disable color output (default to TTY based behavior)

## Building

If you wish to inspect **jlog** and build it by yourself, you may do so
by cloning [this repository](https://github.com/pyr/jlog) and
peforming the following steps :

    mkdir -p $(GOPATH)/src
    cd $(GOPATH)/src && git clone https://github.com/pyr/jlog
    make

### Updating

It uses [godep](https://github.com/golang/dep), so it should be easy.

    dep status
    dep ensure -update

### Example usage

    jlog -f /var/log/some.json.log
