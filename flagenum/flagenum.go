package flagenum

import (
	"cmp"
	"flag"
	"fmt"
	"strings"
)

// CommandLine is the default wrapper of the flag.CommandLine flags.
var CommandLine = New(flag.CommandLine)

// New returns wrapped flagSet.
func New(flagSet *flag.FlagSet) *FlagSetExtension {
	return &FlagSetExtension{flagSet}
}

// Strings defines a string slice flag with specified name, default values, allowed values and usage string.
// The allowed values restrict possible values of the flag.
// If allowed values are defined then an unexpected flag value will cause a panic.
// The return value is the address of a slice that stores values of the flag.
func Strings(name string, defaulValues, allowedValues []string, usage string) *[]string {
	return CommandLine.Strings(name, defaulValues, allowedValues, usage)
}

// String defines a string flag with specified name, default value, allowed values and usage string.
// The allowed values restrict possible value of the flag.
// If allowed values are defined then an unexpected flag value will cause a panic.
// The return value is the address of a slice that stores value of the flag.
func String(name string, value string, allowedValues []string, usage string) *string {
	return CommandLine.String(name, value, allowedValues, usage)
}

// FlagSetExtension extends a FlagSet by addition flag types.
type FlagSetExtension struct {
	*flag.FlagSet
}

// Strings defines a string slice flag with specified name, default values, allowed values and usage string.
// The allowed values restrict possible values of the flag.
// If allowed values are defined then an unexpected flag value will cause a panic.
// The return value is the address of a slice that stores values of the flag.
func (f *FlagSetExtension) Strings(name string, defaulValues, allowedValues []string, usage string) *[]string {
	v, err := Multiple(f.FlagSet, name, defaulValues, allowedValues, strAsIs, strAsIs, usage)
	if err != nil {
		panic(err)
	}
	return v
}

// String defines a string flag with specified name, default value, allowed values and usage string.
// The allowed values restrict possible value of the flag.
// If allowed values are defined then an unexpected flag value will cause a panic.
// The return value is the address of a slice that stores value of the flag.
func (f *FlagSetExtension) String(name string, value string, allowedValues []string, usage string) *string {
	v, err := Single(f.FlagSet, name, value, allowedValues, strAsIs, strAsIs, usage)
	if err != nil {
		panic(err)
	}
	return v
}

func strAsIs(s string) string { return s }

// Value a flag value.
type Value interface {
	cmp.Ordered
}

// Multiple defines a generic slice flag with specified name, default values, allowed values, string converters and usage string.
// The allowed values restrict possible values of the flag.
// Returns the address of a slice that stores values of the flag and an error if something wrong.
func Multiple[V Value](flagSet *flag.FlagSet, name string, defaultValues, allowedValues []V, toVConv func(string) V, toStrConv func(V) string, usage string) (*[]V, error) {
	var result []V
	return &result, MultipleVar(flagSet, &result, name, defaultValues, allowedValues, toVConv, toStrConv, usage)
}

// MultipleVar defines a generic slice flag with specified name, default values, allowed values, string converters and usage string.
// The allowed values restrict possible values of the flag.
// The argument p points to a slice variable in which to store values of the flag.
// Returns an error if something wrong.
func MultipleVar[V Value](flagSet *flag.FlagSet, p *[]V, name string, defaultValues, allowedValues []V, toVConv func(string) V, toStrConv func(V) string, usage string) error {
	allowedUniques, err := getUniques("allowed", name, allowedValues...)
	if err != nil {
		return err
	}
	_, err = getUniques("default", name, defaultValues...)
	if err != nil {
		return err
	}
	if len(allowedValues) > 0 {
		for _, defaultValue := range defaultValues {
			if err := checkDefault(name, defaultValue, allowedValues, allowedUniques, toStrConv); err != nil {
				return err
			}
		}
	}

	*p = append(*p, defaultValues...)
	values := multipleValues[V]{
		name: name, values: p, allowed: allowedValues, uniques: map[V]struct{}{},
		defaults: defaultValues, allowedUniques: allowedUniques, toVConv: toVConv, toStrConv: toStrConv,
	}
	flagSet.Var(&values, name, usage+getSuffix(usage, "any of", toStrConv, allowedValues...))
	return nil
}

// Single defines a generic flag with specified name, default value, allowed values, string converters and usage string.
// The allowed values restrict possible value of the flag.
// Returns the address of a variable that stores value of the flag and an error if something wrong.
func Single[V Value](flagSet *flag.FlagSet, name string, value V, allowedValues []V, toVConv func(string) V, toStrConv func(V) string, usage string) (*V, error) {
	result := value
	return &result, SingleVar(flagSet, &result, name, value, allowedValues, toVConv, toStrConv, usage)
}

// SingleVar defines a generic flag with specified name, default value, allowed values, string converters and usage string.
// The allowed values restrict possible value of the flag.
// The argument p points to a string variable in which to store the value of the flag.
// Returns an error if something wrong.
func SingleVar[V Value](flagSet *flag.FlagSet, p *V, name string, value V, allowedValues []V, toVConv func(string) V, toStrConv func(V) string, usage string) error {
	allowedUniques, err := getUniques("allowed", name, allowedValues...)
	if err != nil {
		return err
	}
	var zero V
	if zero != value {
		if err := checkDefault(name, value, allowedValues, allowedUniques, toStrConv); err != nil {
			return err
		}
	}
	*p = value
	values := singleValue[V]{name: name, value: p, allowed: allowedValues, allowedUniques: allowedUniques, toVConv: toVConv, toStrConv: toStrConv}
	flagSet.Var(&values, name, usage+getSuffix(usage, "one of", toStrConv, allowedValues...))
	return nil
}

func getSuffix[T any](usage, countStr string, toStrConv func(T) string, allowedValues ...T) string {
	suffix := ""
	if len(allowedValues) > 0 {
		suffix = "(allowed `" + countStr + " " + joinToString(toStrConv, allowedValues...) + "`)"
	}
	if len(usage) > 0 {
		suffix = " " + suffix
	}
	return suffix
}

func getUniques[T Value](valueType, name string, values ...T) (map[T]struct{}, error) {
	uniques := map[T]struct{}{}
	for _, e := range values {
		if err := populateUniques(valueType, e, uniques, name); err != nil {
			return uniques, err
		}
	}
	return uniques, nil
}

func joinToString[T any](toStrConv func(T) string, values ...T) string {
	str := strings.Builder{}
	for _, v := range values {
		if str.Len() > 0 {
			str.WriteString(",")
		}
		str.WriteString(toStrConv(v))
	}
	return str.String()
}

func checkDefault[TS ~[]T, T Value](name string, defaultValue T, allowed TS, uniques map[T]struct{}, toStrConv func(T) string) error {
	if err := checkAllowed(defaultValue, allowed, uniques, toStrConv); err != nil {
		return fmt.Errorf("unexpected default value \"%v\" for flag -%s: %w", defaultValue, name, err)
	}
	return nil
}

func checkAllowed[TS ~[]T, T Value](value T, allowed TS, uniques map[T]struct{}, toStrConv func(T) string) error {
	if len(allowed) > 0 {
		if _, ok := uniques[value]; !ok {
			return fmt.Errorf("must be one of %s", joinToString(toStrConv, allowed...))
		}
	}
	return nil
}

func populateUniques[T Value](valueType string, value T, duplicateControl map[T]struct{}, name string) error {
	if _, ok := duplicateControl[value]; !ok {
		duplicateControl[value] = void
		return nil
	}
	if len(valueType) > 0 {
		valueType += " "
	}
	return fmt.Errorf("duplicated %svalue \"%v\" for flag -%s", valueType, value, name)
}

var void struct{}

type multipleValues[T Value] struct {
	name           string
	values         *[]T
	defaults       []T
	allowed        []T
	uniques        map[T]struct{}
	allowedUniques map[T]struct{}
	defaultCleared bool
	toVConv        func(string) T
	toStrConv      func(T) string
}

var _ flag.Value = (*multipleValues[string])(nil)

func (f *multipleValues[T]) String() string {
	v := f.Values()
	c := f.toStrConv
	if v != nil && c != nil {
		return joinToString(f.toStrConv, v...)
	}
	return ""
}

func (f *multipleValues[T]) Set(s string) error {
	if !f.defaultCleared {
		*f.values = nil
		f.defaultCleared = true
	}
	v := f.toVConv(s)
	if err := populateUniques("", v, f.uniques, f.name); err != nil {
		return err
	}
	if err := checkAllowed(v, f.allowed, f.allowedUniques, f.toStrConv); err != nil {
		return err
	}
	*f.values = append(*f.values, v)
	return nil
}

func (f *multipleValues[T]) Get() any {
	return f.Values()
}

func (f *multipleValues[T]) Values() []T {
	if v := f.values; v != nil && len(*v) > 0 {
		return *v
	}
	return f.defaults
}

type singleValue[T Value] struct {
	name           string
	value          *T
	allowed        []T
	allowedUniques map[T]struct{}
	toVConv        func(string) T
	toStrConv      func(T) string
}

var _ flag.Value = (*singleValue[string])(nil)

func (f *singleValue[T]) String() string {
	v := f.Value()
	c := f.toStrConv
	if v != nil && c != nil {
		return joinToString(c, *v)
	}
	return ""
}

func (f *singleValue[T]) Set(s string) error {
	v := f.toVConv(s)
	if err := checkAllowed(v, f.allowed, f.allowedUniques, f.toStrConv); err != nil {
		return err
	}
	*f.value = v
	return nil
}

func (f *singleValue[T]) Get() any {
	return f.Value()
}

func (f *singleValue[T]) Value() *T {
	return f.value
}
