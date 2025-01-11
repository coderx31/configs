package configs

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
)

const SensitiveDataMaskString = "***************"

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

		secretTag, ok := structField.Tag.Lookup("secret")
		if ok && secretTag == "true" {
			data = append(data, []string{structField.Name, SensitiveDataMaskString})
		} else {
			data = append(data, []string{structField.Name, field.String()})
		}
	}

	table.SetHeader([]string{"Config", "Value"})
	table.AppendBulk(data)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Render()
}
