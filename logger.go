package main

import (
	"log"
	"os"
)

const ResetColour string = "\033[0m"
const RedColour string = "\033[31m"
const GreenColour string = "\033[32m"
const YellowColour string = "\033[33m"
const BlueColour string = "\033[34m"
const PurpleColour string = "\033[35m"
const CyanColour string = "\033[36m"
const WhiteColour string = "\033[37m"

type logger struct {
	*log.Logger
	debug bool
}

var Log logger

func (l *logger) Ok(m string) {
	l.Println(GreenColour + m + ResetColour)
}

func (l *logger) Err(m string) {
	l.Fatal(RedColour + m + ResetColour)
}

func (l *logger) Warn(m string) {
	l.Println(YellowColour + m + ResetColour)
}


func (l *logger) Debug(m string) {
	if l.debug {
		l.Println(WhiteColour + m + ResetColour)
	}
}

func (l *logger) DebugOn() {
	l.debug = true
}

func (l *logger) DebugOff() {
	l.debug = false
}

func init() {
	Log.Logger = log.New(os.Stderr, "ss-cloud-client: ", log.Ldate|log.Ltime|log.Lmsgprefix)
}
