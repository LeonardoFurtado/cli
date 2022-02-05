package shared

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var rangeRE = regexp.MustCompile(`(>|>=|<|<=|\*\.\.)?\d+(\.\.\*|\.\.\d+)?`)

type validator func(string) error

type qualifier struct {
	key       string
	kind      string
	set       bool
	validator validator
	value     string
}

type parameter = qualifier

func NewQualifier(key, kind, value string, validator func(string) error) *qualifier {
	return &qualifier{
		key:       key,
		kind:      kind,
		validator: validator,
		value:     value,
	}
}

func NewParameter(key, kind, value string, validator func(string) error) *parameter {
	return &parameter{
		key:       key,
		kind:      kind,
		validator: validator,
		value:     value,
	}
}

func (q *qualifier) IsSet() bool {
	return q.set
}

func (q *qualifier) Key() string {
	return q.key
}

func (q *qualifier) Set(v string) error {
	if q.validator != nil {
		err := q.validator(v)
		if err != nil {
			return err
		}
	}
	q.set = true
	q.value = v
	return nil
}

func (q *qualifier) String() string {
	return q.value
}

func (q *qualifier) Type() string {
	return q.kind
}

// Validate that value is one of a list of options
func OptsValidator(opts []string) validator {
	return func(v string) error {
		if !isIncluded(v, opts) {
			return fmt.Errorf("%s is not included in %s", v, strings.Join(opts, ", "))
		}
		return nil
	}
}

// Validate that each value in v matches a value in list of options
func MultiOptsValidator(opts []string) validator {
	return func(v string) error {
		s := strings.Split(v, ",")
		for _, t := range s {
			if !isIncluded(t, opts) {
				return fmt.Errorf("%q is not included in %s", t, strings.Join(opts, ", "))
			}
		}
		return nil
	}
}

func isIncluded(v string, opts []string) bool {
	for _, opt := range opts {
		if v == opt {
			return true
		}
	}
	return false
}

// Validate that value is a boolean
func BoolValidator() validator {
	return func(v string) error {
		_, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("%s is not a boolean value", v)
		}
		return nil
	}
}

// Validate that value is a correct range format
func RangeValidator() validator {
	return func(v string) error {
		if !rangeRE.MatchString(v) {
			return fmt.Errorf("%s is invalid format", v)
		}
		return nil
	}
}

// Validate that value is a correct date format
//TODO: write regex here
func DateValidator() validator {
	return func(v string) error {
		if !rangeRE.MatchString(v) {
			return fmt.Errorf("%s is invalid format", v)
		}
		return nil
	}
}
