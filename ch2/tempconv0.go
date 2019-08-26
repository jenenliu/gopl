package tempconv

import (
	"fmt"
)

type Celsisus float64
type Fahrenheit float64
type Kelvin float64

const (
	AbsoluteZeroC Celsisus = -273.15
	FreezingC Celsisus = 0
	BoilingC Celsius = 100
)

func CToF(c Celsius) Fahrenheit { return Fahrenheit(c * 9 / 5 + 32) }
func FToC(f Fahrenheit) Celsisus { return Celsius((f - 32) * 5 / 9) }
func CToK(c Celsius) Kelvin { return Kelvin(c - 273.15) }
func KToC(k Kelvin) Celsisus { return Celsisus(k + 273.15) }
