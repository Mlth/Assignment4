#!/bin/sh
xterm -e "go run peer.go 0 5" &
xterm -e "go run peer.go 1 5" &
xterm -e "go run peer.go 2 5" &
xterm -e "go run peer.go 3 5" &
xterm -e "go run peer.go 4 5" &