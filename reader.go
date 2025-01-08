package configs

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
)

const (
	defaultLength    = 10
	defaultAddLength = 5
)

type Reader interface {
	Register() error
}

type Validator interface {
	Validation() error
}

type Printer interface {
	Print() interface{}
}

func Load(configs ...Reader) error {
	for _, config := range configs {
		err := config.Register()
		if err != nil {
			return err
		}

		validation, ok := config.(Validator)
		if ok {
			err = validation.Validation()
			if err != nil {
				return err
			}
		}

		printer, ok := config.(Printer)
		if ok {
			printTable(printer)
		}
	}
	return nil
}

func printTable(p Printer) {
	table := tablewriter.NewWriter(os.Stdout)

	var data [][]string

	printer := p.Print()

	values := reflect.ValueOf(printer)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}

	if values.Kind() == reflect.Interface {
		values = values.Elem()
	}

	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		structField := values.Type().Field(i)

		_, ok := structField.Tag.Lookup("secret")
		if ok {
			data = append(data, []string{structField.Name, mask(field.String())})
		} else {
			data = append(data, []string{structField.Name, field.String()})
		}
	}

	table.SetHeader([]string{"Config", "Value"})
	table.AppendBulk(data)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Render()
}

func mask(value string) string {
	length := defaultLength
	if len(value) > defaultLength {
		length = len(value) + defaultAddLength
	}
	runes := make([]rune, length)

	for i := 0; i < length; i++ {
		runes[i] = '*'
	}

	return string(runes)
}
