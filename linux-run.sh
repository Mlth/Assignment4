#!/bin/sh
gnome-terminal -e "go run peer.go 0 5" &
gnome-terminal -e "go run peer.go 1 5" &
gnome-terminal -e "go run peer.go 2 5" &
gnome-terminal -e "go run peer.go 3 5" &
gnome-terminal -e "go run peer.go 4 5" &