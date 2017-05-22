package main

type Options struct {
	target  string
	verbose bool
}

func NewOptions() *Options {
	o := new(Options)
	o.target = ""
	o.verbose = false
	return o
}

func (o *Options) Parse(arguments []string) {
	for _, v := range arguments {
		switch v {
		case "-v":
			o.verbose = true
		default:
			o.target = v
		}
	}
}
