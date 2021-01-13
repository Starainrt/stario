package stario

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type InputMsg struct {
	msg string
	err error
}

func MessageBox(hint string, defaultVal string) InputMsg {
	if hint != "" {
		fmt.Print(hint)
	}
	inputReader := bufio.NewReader(os.Stdin)
	str, err := inputReader.ReadString('\n')
	if err != nil {
		return InputMsg{"", err}
	}
	return InputMsg{strings.TrimSpace(str), nil}
}

func (im InputMsg) String() (string, error) {
	if im.err != nil {
		return "", im.err
	}
	return im.msg, nil
}

func (im InputMsg) MustString() string {
	res, _ := im.String()
	return res
}

func (im InputMsg) Int() (int, error) {
	if im.err != nil {
		return 0, im.err
	}
	return strconv.Atoi(im.msg)
}

func (im InputMsg) MustInt() int {
	res, _ := im.Int()
	return res
}

func (im InputMsg) Int64() (int64, error) {
	if im.err != nil {
		return 0, im.err
	}
	return strconv.ParseInt(im.msg, 10, 64)
}

func (im InputMsg) MustInt64() int64 {
	res, _ := im.Int64()
	return res
}

func (im InputMsg) Uint64() (uint64, error) {
	if im.err != nil {
		return 0, im.err
	}
	return strconv.ParseUint(im.msg, 10, 64)
}

func (im InputMsg) MustUint64() uint64 {
	res, _ := im.Uint64()
	return res
}

func (im InputMsg) Bool() (bool, error) {
	if im.err != nil {
		return false, im.err
	}
	return strconv.ParseBool(im.msg)
}

func (im InputMsg) MustBool() bool {
	res, _ := im.Bool()
	return res
}

func (im InputMsg) Float64() (float64, error) {
	if im.err != nil {
		return 0, im.err
	}
	return strconv.ParseFloat(im.msg, 64)
}

func (im InputMsg) MustFloat64() float64 {
	res, _ := im.Float64()
	return res
}

func (im InputMsg) Float32() (float32, error) {
	if im.err != nil {
		return 0, im.err
	}
	res, err := strconv.ParseFloat(im.msg, 32)
	return float32(res), err
}

func (im InputMsg) MustFloat32() float32 {
	res, _ := im.Float32()
	return res
}

func YesNo(hint string, defaults bool) bool {
	res := strings.ToUpper(MessageBox(hint, "").MustString())
	if res == "" {
		return defaults
	}
	res = res[0:1]
	if res == "Y" {
		return true
	} else if res == "N" {
		return false
	} else {
		return defaults
	}
}

func StopUntil(hint string, trigger string) error {
	pressLen := len(trigger)
	if trigger == "" {
		pressLen = 1
	} else {
		pressLen = len(trigger)
	}
	ioBuf := make([]byte, pressLen)
	for {
		if hint != "" {
			fmt.Print(hint)
		}
		inputReader := bufio.NewReader(os.Stdin)
		_, err := inputReader.Read(ioBuf)
		if err != nil {
			return err
		}
		if string(ioBuf) == trigger || trigger == "" {
			break
		}
		fmt.Print("\n")
	}
	return nil
}
