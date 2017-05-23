package main

type Options struct {
	clear   bool
	target  string
	verbose bool
}

func NewOptions() *Options {
	o := new(Options)
	o.clear = false
	o.target = ""
	o.verbose = false
	return o
}

func (o *Options) Parse(arguments []string) {
	for _, v := range arguments {
		switch v {
		case "-c":
			o.clear = true
		case "-v":
			o.verbose = true
		default:
			o.target = v
		}
	}
}
