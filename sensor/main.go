package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/ina228"
)

func main() {
	machine.I2C0.Configure(machine.I2CConfig{})

	dev := ina228.New(machine.I2C0, 0)
	dev.Configure(ina228.Config{
		ADCRange:            1, // 4x precision
		MaxCurrentMilliAmps: 2000,
	})

	if !dev.Connected() {
		println("INA228 not found")
		return
	}

	println("Time ms, Current mA, Power mW, Energy mJ, Charge mC, Temp C")

	for {
		shuntnv := dev.ShuntVoltage()
		busnv := dev.BusVoltage()
		currentnv := dev.Current()
		powernv := dev.Power()
		energynj := dev.Energy()
		chargenc := dev.Charge()
		tempuc := dev.Temperature()

		println(fmtD(time.Now().UnixNano()/1000, 9, 3) + "," +
			fmtD(shuntnv, 4, 6) + "," +
			fmtD(busnv, 7, 6) + "," +
			fmtD(currentnv, 4, 6) + "," +
			fmtD(powernv, 4, 6) + "," +
			fmtD(energynj, 9, 6) + "," +
			fmtD(chargenc, 9, 6) + "," +
			fmtD(int64(tempuc), 3, 6))

		time.Sleep(100 * time.Millisecond)
	}
}

func fmtD(val int64, i int, f int) string {
	result := make([]byte, i+f+1)
	neg := false

	if val < 0 {
		val = -val
		neg = true
	}

	for p := len(result) - 1; p >= 0; p-- {
		result[p] = byte(int32('0') + int32(val%10))
		val = val / 10

		if p == i+1 && p > 0 {
			p--
			result[p] = '.'
		}
	}

	if neg {
		result[0] = '-'
	}

	return string(result)
}
