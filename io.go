package stario

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

type InputMsg struct {
	msg string
	err error
}

func Passwd(hint string, defaultVal string) InputMsg {
	return passwd(hint, defaultVal, "")
}

func PasswdWithMask(hint string, defaultVal string, mask string) InputMsg {
	return passwd(hint, defaultVal, mask)
}

func MessageBoxRaw(hint string, defaultVal string) InputMsg {
	return messageBox(hint, defaultVal)
}

func messageBox(hint string, defaultVal string) InputMsg {
	var ioBuf []rune
	if hint != "" {
		fmt.Print(hint)
	}
	if strings.Index(hint, "\n") >= 0 {
		hint = strings.TrimSpace(hint[strings.LastIndex(hint, "\n"):])
	}
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return InputMsg{"", err}
	}
	defer fmt.Println()
	defer terminal.Restore(fd, state)
	inputReader := bufio.NewReader(os.Stdin)
	for {
		b, _, err := inputReader.ReadRune()
		if err != nil {
			return InputMsg{"", err}
		}
		if b == 0x0d {
			strValue := strings.TrimSpace(string(ioBuf))
			if len(strValue) == 0 {
				strValue = defaultVal
			}
			return InputMsg{strValue, nil}
		}
		if b == 0x08 || b == 0x7F {
			if len(ioBuf) > 0 {
				ioBuf = ioBuf[:len(ioBuf)-1]
			}
			fmt.Print("\r")
			for i := 0; i < len(ioBuf)+2+len(hint); i++ {
				fmt.Print(" ")
			}
		} else {
			ioBuf = append(ioBuf, b)
		}
		fmt.Print("\r")
		if hint != "" {
			fmt.Print(hint)
		}
		fmt.Print(string(ioBuf))
	}
}

func passwd(hint string, defaultVal string, mask string) InputMsg {
	var ioBuf []rune
	if hint != "" {
		fmt.Print(hint)
	}
	if strings.Index(hint, "\n") >= 0 {
		hint = strings.TrimSpace(hint[strings.LastIndex(hint, "\n"):])
	}
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return InputMsg{"", err}
	}
	defer fmt.Println()
	defer terminal.Restore(fd, state)
	inputReader := bufio.NewReader(os.Stdin)
	for {
		b, _, err := inputReader.ReadRune()
		if err != nil {
			return InputMsg{"", err}
		}
		if b == 0x0d {
			strValue := strings.TrimSpace(string(ioBuf))
			if len(strValue) == 0 {
				strValue = defaultVal
			}
			return InputMsg{strValue, nil}
		}
		if b == 0x08 || b == 0x7F {
			if len(ioBuf) > 0 {
				ioBuf = ioBuf[:len(ioBuf)-1]
			}
			fmt.Print("\r")
			for i := 0; i < len(ioBuf)+2+len(hint); i++ {
				fmt.Print(" ")
			}
		} else {
			ioBuf = append(ioBuf, b)
		}
		fmt.Print("\r")
		if hint != "" {
			fmt.Print(hint)
		}
		for i := 0; i < len(ioBuf); i++ {
			fmt.Print(mask)
		}
	}
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
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		str = defaultVal
	}
	return InputMsg{str, nil}
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
	for {
		res := strings.ToUpper(MessageBox(hint, "").MustString())
		if res == "" {
			return defaults
		}
		res = res[0:1]
		if res == "Y" {
			return true
		} else if res == "N" {
			return false
		}
	}
}

func StopUntil(hint string, trigger string, repeat bool) error {
	pressLen := len([]rune(trigger))
	if trigger == "" {
		pressLen = 1
	}
	fd := int(os.Stdin.Fd())
	if hint != "" {
		fmt.Print(hint)
	}
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, state)
	inputReader := bufio.NewReader(os.Stdin)
	//ioBuf := make([]byte, pressLen)
	i := 0
	for {
		b, _, err := inputReader.ReadRune()
		if err != nil {
			return err
		}
		if trigger == "" {
			break
		}
		if b == []rune(trigger)[i] {
			i++
			if i == pressLen {
				break
			}
			continue
		}
		i = 0
		if hint != "" && repeat {
			fmt.Print("\r\n")
			fmt.Print(hint)
		}
	}
	return nil
}
