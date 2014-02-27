package log

import (
	//lg "code.google.com/p/log4go"
	"log"
)

func init() {
	//lg.LoadConfiguration("log.xml")
}

// Tracef formats message according to format specifier
// and writes to default logger with log level = Trace.
func Tracef(format string, params ...interface{}) {
	//lg.Trace(format, params...)
	log.Printf(format, params...)
}

// Debugf formats message according to format specifier
// and writes to default logger with log level = Debug.
func Debugf(format string, params ...interface{}) {
	log.Printf(format, params...)
}

// Infof formats message according to format specifier
// and writes to default logger with log level = Info.
func Infof(format string, params ...interface{}) {
	log.Printf(format, params...)
}

// Warnf formats message according to format specifier and writes to default logger with log level = Warn
func Warnf(format string, params ...interface{}) {
	log.Printf(format, params...)
}

// Errorf formats message according to format specifier and writes to default logger with log level = Error
func Errorf(format string, params ...interface{}) {
	log.Printf(format, params...)
}

// Trace formats message using the default formats for its operands and writes to default logger with log level = Trace
func Trace(v ...interface{}) {
	log.Print(v)
}

// Debug formats message using the default formats for its operands and writes to default logger with log level = Debug
func Debug(v ...interface{}) {
	log.Print(v)
}

// Info formats message using the default formats for its operands and writes to default logger with log level = Info
func Info(v ...interface{}) {
	log.Print(v)
}
func Infoln(v ...interface{}) {
	log.Println(v)
}

// Warn formats message using the default formats for its operands and writes to default logger with log level = Warn
func Warn(v ...interface{}) {
	log.Print(v)
}

// Error formats message using the default formats for its operands and writes to default logger with log level = Error
func Error(v ...interface{}) {
	log.Print(v)
}
