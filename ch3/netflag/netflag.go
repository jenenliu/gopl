package main

import "fmt"

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Flags uint

const (
	FlagUp Flags = 1 << iota // is up
	FlagBroadcast
	FlagLoopback
	FlagPointToPoint
	FlagMulticast
)

func main() {
	fmt.Printf("%d %d %d %d %d %d %d\n", Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday)
	fmt.Printf("%b %b %b %b %b\n", FlagUp, FlagBroadcast, FlagLoopback, FlagPointToPoint, FlagMulticast)
}
