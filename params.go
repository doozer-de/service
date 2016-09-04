package service

import (
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Params struct {
	set *flag.FlagSet
}

func (p *Params) Parse(args url.Values) error {
	if p.set == nil {
		p.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}

FLAG_LOOP:
	for name, vals := range args {
		for _, v := range vals {

			f := p.set.Lookup(name)
			if f == nil {
				continue FLAG_LOOP
			}

			// Check if the value is empty
			if v == "" {
				if bv, ok := f.Value.(boolFlag); ok && bv.IsBoolFlag() {
					bv.Set("true")

					continue FLAG_LOOP
				}
			}

			err := p.set.Set(name, v)
			if err != nil {
				if !strings.Contains(err.Error(), "no such flag -") {
					err = fmt.Errorf("bad param '%s': %s", name, err.Error())
					return err
				}
			}
		}
	}

	return nil
}

func (p *Params) Usage() map[string][3]string {
	docs := make(map[string][3]string)
	var translations map[string]string = map[string]string{
		"*flag.stringValue":   "string",
		"*flag.durationValue": "duration",
		"*flag.intValue":      "int",
		"*flag.boolValue":     "bool",
		"*flag.float64Value":  "float64",
		"*flag.int64Value":    "int64",
		"*flag.uintValue":     "uint",
		"*flag.uint64Value":   "uint64",
		"*service.SString":    "[]string",
		"*service.SDuration":  "[]duration",
		"*service.SInt":       "[]int",
		"*service.SBool":      "[]bool",
		"*service.SFloat64":   "[]float64",
		"*service.SInt64":     "[]int64",
		"*service.SUint":      "[]uint",
		"*service.SUint64":    "[]uint64",
	}
	p.set.VisitAll(func(flag *flag.Flag) {
		niceName := translations[fmt.Sprintf("%T", flag.Value)]
		if niceName == "" {
			niceName = fmt.Sprintf("%T", flag.Value)
		}
		docs[flag.Name] = [...]string{flag.Name, niceName, flag.Usage}
	})
	return docs
}

type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

func (p *Params) Bool(name string, value bool, usage string) *bool {
	if p.set == nil {
		p.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	b := new(bool)
	p.set.BoolVar(b, name, value, usage)
	return b
}

type SBool []bool

func (s *SBool) String() string {
	return fmt.Sprint(*s)
}

func (s *SBool) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseBool(dt)
		if err != nil {
			return err
		}
		*s = append(*s, parsed)
	}
	return nil
}

func (p *Params) SBool(name string, value bool, usage string) *SBool {
	if p.set == nil {
		p.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	b := new(SBool)
	p.set.Var(b, name, usage)
	return b
}

func (rp *Params) Int(name string, value int, usage string) *int {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(int)
	rp.set.IntVar(p, name, value, usage)
	return p
}

type SInt []int

func (s *SInt) String() string {
	return fmt.Sprint(*s)
}

func (s *SInt) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseInt(dt, 0, 64)
		if err != nil {
			return err
		}
		*s = append(*s, int(parsed))
	}
	return nil
}

func (rp *Params) SInt(name string, value int, usage string) *SInt {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SInt)
	rp.set.Var(p, name, usage)
	return p
}

func (rp *Params) Int64(name string, value int64, usage string) *int64 {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(int64)
	rp.set.Int64Var(p, name, value, usage)
	return p
}

type SInt64 []int64

func (s *SInt64) String() string {
	return fmt.Sprint(*s)
}

func (s *SInt64) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseInt(dt, 0, 64)
		if err != nil {
			return err
		}
		*s = append(*s, int64(parsed))
	}
	return nil
}

func (rp *Params) SInt64(name string, value int64, usage string) *SInt64 {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SInt64)
	rp.set.Var(p, name, usage)
	return p
}

func (rp *Params) Uint(name string, value uint, usage string) *uint {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(uint)
	rp.set.UintVar(p, name, value, usage)
	return p
}

type SUint []uint

func (s *SUint) String() string {
	return fmt.Sprint(*s)
}

func (s *SUint) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseFloat(dt, 64)
		if err != nil {
			return err
		}
		*s = append(*s, uint(parsed))
	}
	return nil
}

func (rp *Params) SUint(name string, value uint, usage string) *SUint {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SUint)
	rp.set.Var(p, name, usage)
	return p
}

func (rp *Params) Uint64(name string, value uint64, usage string) *uint64 {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(uint64)
	rp.set.Uint64Var(p, name, value, usage)
	return p
}

type SUint64 []uint64

func (s *SUint64) String() string {
	return fmt.Sprint(*s)
}

func (s *SUint64) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseFloat(dt, 64)
		if err != nil {
			return err
		}
		*s = append(*s, uint64(parsed))
	}
	return nil
}

func (rp *Params) SUint64(name string, value uint64, usage string) *SUint64 {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SUint64)
	rp.set.Var(p, name, usage)
	return p
}

func (rp *Params) String(name string, value string, usage string) *string {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(string)
	rp.set.StringVar(p, name, value, usage)
	return p
}

type SString []string

func (s *SString) String() string {
	return strings.Join(*s, ",")
}

func (s *SString) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		*s = append(*s, dt)
	}
	return nil
}

func (rp *Params) SString(name string, value string, usage string) *SString {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SString)
	rp.set.Var(p, name, usage)
	return p
}

func (rp *Params) Float64(name string, value float64, usage string) *float64 {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(float64)
	rp.set.Float64Var(p, name, value, usage)
	return p
}

type SFloat64 []float64

func (s *SFloat64) String() string {
	return fmt.Sprintf("%f", *s)
}

func (s *SFloat64) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseFloat(dt, 64)
		if err != nil {
			return err
		}
		*s = append(*s, parsed)
	}
	return nil
}

func (rp *Params) SFloat64(name string, value float64, usage string) *SFloat64 {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SFloat64)
	rp.set.Var(p, name, usage)
	return p
}

func (rp *Params) Duration(name string, value time.Duration, usage string) *time.Duration {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(time.Duration)
	rp.set.DurationVar(p, name, value, usage)
	return p
}

type SDuration []time.Duration

func (s *SDuration) String() string {
	return fmt.Sprint(*s)
}

func (s *SDuration) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := time.ParseDuration(dt)
		if err != nil {
			return err
		}
		*s = append(*s, parsed)
	}
	return nil
}

func (rp *Params) SDuration(name string, value time.Duration, usage string) *SDuration {
	if rp.set == nil {
		rp.set = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SDuration)
	rp.set.Var(p, name, usage)
	return p
}
