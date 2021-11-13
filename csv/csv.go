package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func ToCsvData(path string, csvData ...CsvData) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("open file is failed, err: ", err)
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			return
		}
	}()
	// 写入UTF-8 BOM，防止中文乱码
	if _, err = file.WriteString("\xEF\xBB\xBF"); err != nil {
		return err
	}
	w := csv.NewWriter(file)
	for _, data := range csvData {
		if err = w.WriteAll(data.ToCsvData()); err != nil {
			//write failed do something
			return err
		}
	}
	w.Flush()
	return nil
}

type CsvData []interface{}

func (c CsvData) ToCsvData() (ds [][]string) {
	for _, data := range c {
		typeOf := reflect.TypeOf(data)
		value := reflect.ValueOf(data)
		switch typeOf.Kind() {
		case reflect.Struct:
			var keys []string
			num := typeOf.NumField()
			for i := 0; i < num; i++ {
				key := typeOf.Field(i).Name
				keys = append(keys, key)
			}
			ds = append(ds, keys)
			var values []string
			num = value.NumField()
			for i := 0; i < num; i++ {
				values = append(values, c.getValue(value.Field(i))...)
			}
			ds = append(ds, values)
		case reflect.Slice:
			var values []string
			num := value.Len()
			for i := 0; i < num; i++ {
				values = append(values, c.getValue(value.Index(i))...)
			}
			ds = append(ds, values)
		case reflect.Map:
			var keys []string
			var values []string
			iter := value.MapRange()
			for iter.Next() {
				keys = append(keys, c.getValue(iter.Key())...)
				values = append(values, c.getValue(iter.Value())...)
			}
			ds = append(ds, keys, values)
		}
	}
	return
}

func (c CsvData) getValue(value reflect.Value) (values []string) {
	switch value.Kind() {
	case reflect.String:
		values = append(values, value.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		values = append(values, strconv.FormatInt(value.Int(), 10))
	case reflect.Float32, reflect.Float64:
		values = append(values, strconv.FormatFloat(value.Float(), 'f', 2, 64))
	case reflect.Interface:
		values = append(values, c.getValue(value.Elem())...)
	default: // 其他类型暂不支持

	}
	return
}
