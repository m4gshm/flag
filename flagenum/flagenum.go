package flagenum

import (
	"cmp"
	"flag"
	"fmt"
	"slices"
	"strings"
)

var CommandLine = New(flag.CommandLine)

func New(flagSet *flag.FlagSet) *FlagSetExtension {
	return &FlagSetExtension{flagSet}
}

func MultipleStrings(name string, usage string) *[]string {
	return CommandLine.MultipleStringsDefStrict(name, usage, nil, nil)
}

func MultiStringsDef(name string, usage string, defaultValues []string) *[]string {
	return CommandLine.MultipleStringsDefStrict(name, usage, defaultValues, nil)
}

func MultiStringsStrict(name string, usage string, allowedValues []string) *[]string {
	return CommandLine.MultipleStringsDefStrict(name, usage, nil, allowedValues)
}

func MultipleStringsDefStrict(name string, usage string, defaulValues, allowedValues []string) *[]string {
	return CommandLine.MultipleStringsDefStrict(name, usage, defaulValues, allowedValues)
}

type Value interface {
	cmp.Ordered
}

type FlagSetExtension struct {
	*flag.FlagSet
}

func (f *FlagSetExtension) MultipleStringsDefStrict(name string, usage string, defaulValues, allowedValues []string) *[]string {
	v, err := Multiple(f.FlagSet, name, usage, func(s string) string { return s }, defaulValues, allowedValues)
	if err != nil {
		panic(err)
	}
	return v
}

func Multiple[V Value](f *flag.FlagSet, name string, usage string, converter func(string) V, defaultValues, allowedValues []V) (*[]V, error) {
	allowed, err := getUniques("allowed", name, allowedValues...)
	if err != nil {
		return nil, err
	}
	_, err = getUniques("default", name, defaultValues...)
	if err != nil {
		return nil, err
	}
	if len(allowedValues) > 0 {
		for _, defaultValue := range defaultValues {
			if err := checkAllowed(defaultValue, allowedValues, allowed); err != nil {
				return nil, fmt.Errorf("unexpected default value \"%v\" for flag -%s: %w", defaultValue, name, err)
			}
		}
	}
	values := multipleValues[V]{
		name: name, values: slices.Clone(defaultValues), allowed: allowedValues, uniques: map[V]struct{}{},
		defaults: defaultValues, allowedUniques: allowed, converter: converter,
	}
	f.Var(&values, name, usage+getSuffix(usage, allowedValues...))
	return &values.values, nil
}

func Single[V Value](f *flag.FlagSet, name string, usage string, converter func(string) V, allowedValues []V) (*V, error) {
	allowed, err := getUniques("allowed", name, allowedValues...)
	if err != nil {
		return nil, err
	}
	var zero V
	values := singleValue[V]{name: name, value: zero, allowed: allowedValues, allowedUniques: allowed, converter: converter}
	f.Var(&values, name, usage+getSuffix(usage, allowedValues...))
	return &values.value, nil
}

func getSuffix[V Value](usage string, allowedValues ...V) string {
	suffix := ""
	if len(allowedValues) > 0 {
		suffix = "(allowed: " + joinToString(allowedValues) + ")"
	}
	if len(usage) > 0 {
		suffix = " " + suffix
	}
	return suffix
}

func getUniques[V Value](valueType, name string, values ...V) (map[V]struct{}, error) {
	uniques := map[V]struct{}{}
	for _, e := range values {
		if err := populateUniques(valueType, e, uniques, name); err != nil {
			return uniques, err
		}
	}
	return uniques, nil
}

func joinToString[T any](values ...T) string {
	str := strings.Builder{}
	for _, v := range values {
		if str.Len() > 0 {
			str.WriteString(",")
		}
		str.WriteString(fmt.Sprintf("%v", v))
	}
	return str.String()
}

func checkAllowed[TS ~[]T, T Value](value T, allowed TS, uniques map[T]struct{}) error {
	if len(allowed) > 0 {
		if _, ok := uniques[value]; !ok {
			return fmt.Errorf("must be one of %s", joinToString(allowed))
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
	values         []T
	defaults       []T
	allowed        []T
	uniques        map[T]struct{}
	allowedUniques map[T]struct{}
	defaultCleared bool
	converter      func(string) T
}

var _ flag.Value = (*multipleValues[string])(nil)

func (f *multipleValues[T]) String() string {
	return joinToString(f.Values())
}

func (f *multipleValues[T]) Set(s string) error {
	if !f.defaultCleared {
		f.values = nil
		f.defaultCleared = true
	}
	v := f.converter(s)
	if err := populateUniques("", v, f.uniques, f.name); err != nil {
		return err
	}
	if err := checkAllowed(v, f.allowed, f.allowedUniques); err != nil {
		return err
	}
	f.values = append(f.values, v)
	return nil
}

func (f *multipleValues[T]) Get() any {
	return f.Values()
}

func (f *multipleValues[T]) Values() []T {
	if v := f.values; len(v) > 0 {
		return v
	}
	return f.defaults
}

type singleValue[T Value] struct {
	name           string
	value          T
	allowed        []T
	allowedUniques map[T]struct{}
	converter      func(string) T
}

var _ flag.Value = (*multipleValues[string])(nil)

func (f *singleValue[T]) String() string {
	return joinToString(f.Value())
}

func (f *singleValue[T]) Set(s string) error {
	v := f.converter(s)
	if err := checkAllowed(v, f.allowed, f.allowedUniques); err != nil {
		return err
	}
	f.value = v
	return nil
}

func (f *singleValue[T]) Get() any {
	return f.Value()
}

func (f *singleValue[T]) Value() *T {
	return &f.value
}
