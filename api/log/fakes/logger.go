package fakes

import "fmt"

type Logger struct{}

func (*Logger) Debugf(format string, args ...interface{}) {

}

func (*Logger) Infof(format string, args ...interface{}) {

}

func (*Logger) Warnf(format string, args ...interface{}) {

}

func (*Logger) Errorf(format string, args ...interface{}) {

}

func (*Logger) Debug(args ...interface{}) {
	fmt.Printf("\nTEST DEBUG LOG: %v\n", args)
}

func (*Logger) Info(args ...interface{}) {
	fmt.Printf("\nTEST INFO LOG: %v\n", args)
}

func (*Logger) Warn(args ...interface{}) {
	fmt.Printf("\nTEST WARN LOG: %v\n", args)
}

func (*Logger) Error(args ...interface{}) {
	fmt.Printf("\nTEST ERROR LOG: %v\n", args)
}
