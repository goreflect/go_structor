package pipeline

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	gohocon "github.com/goreflect/go_hocon"
	"github.com/goreflect/gostructor/tags"
)

type HoconConfig struct {
	fileName            string
	configureFileParsed *gohocon.Config
}

func (config HoconConfig) getElementName(context *structContext) string {
	currentTagHoconValue := context.StructField.Tag.Get(tags.TagHocon)
	if strings.Contains(context.Prefix, currentTagHoconValue) {
		fmt.Println("[HOCON]: Level: debug. Current field name: ", context.Prefix)
		return context.Prefix
	}
	returnName := context.Prefix + "."
	if currentTagHoconValue == "" {
		returnName += context.StructField.Name
	}

	returnName += currentTagHoconValue
	fmt.Println("[HOCON]: Level: debug. Current field name: ", returnName)
	return returnName
}

func (config HoconConfig) GetComplexType(context *structContext) GoStructorValue {
	loaded, structValue := config.typeSafeLoadConfigFile(context)
	if !loaded {
		return *structValue
	}
	valueIndirect := reflect.Indirect(context.Value)
	switch valueIndirect.Kind() {
	case reflect.Array:
		return config.getSliceFromHocon(context)
	case reflect.Slice:
		return config.getSliceFromHocon(context)
	case reflect.Map:
		return config.getMapFromHocon(context)
	default:
		return config.GetBaseType(context)
	}
}

// return true - if loaded config or successfully load config by filename
func (config *HoconConfig) typeSafeLoadConfigFile(context *structContext) (bool, *GoStructorValue) {
	if config.configureFileParsed == nil {
		configParsed, err := gohocon.LoadConfig(config.fileName)
		if err != nil {
			return false, &GoStructorValue{notAValue: &NotAValue{
				ValueAddress: context.Value.Interface(),
				Error:        err,
			},
			}
		}
		config.configureFileParsed = configParsed
		return true, nil
	}
	return true, nil
}

func (config *HoconConfig) GetBaseType(context *structContext) GoStructorValue {
	configLoad, structValue := config.typeSafeLoadConfigFile(context)
	if !configLoad {
		return *structValue
	}
	valueIndirect := reflect.Indirect(context.Value)
	path := config.getElementName(context)
	switch valueIndirect.Kind() {
	case reflect.Int:
		loadValue, errLoading := config.configureFileParsed.GetString(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		valueParsedFromConfigFile, errParsed := strconv.ParseInt(loadValue, 10, 32)
		if errParsed != nil {
			return NewGoStructorNoValue(valueIndirect, errParsed)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(int(valueParsedFromConfigFile)))
	case reflect.Int8:
		loadValue, errLoading := config.configureFileParsed.GetString(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		valueParsedFromConfigFile, errParsed := strconv.ParseInt(loadValue, 10, 8)
		if errParsed != nil {
			return NewGoStructorNoValue(valueIndirect, errParsed)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(int8(valueParsedFromConfigFile)))
	case reflect.Int16:
		loadValue, errLoading := config.configureFileParsed.GetString(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		valueParsedFromConfigFile, errParsed := strconv.ParseInt(loadValue, 10, 16)
		if errParsed != nil {
			return NewGoStructorNoValue(valueIndirect, errParsed)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(int16(valueParsedFromConfigFile)))
	case reflect.Int32:
		loadValue, errLoading := config.configureFileParsed.GetInt32(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(loadValue))
	case reflect.Int64:
		loadValue, errLoading := config.configureFileParsed.GetInt64(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(loadValue))
	case reflect.Float32:
		loadValue, errLoading := config.configureFileParsed.GetFloat32(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(loadValue))
	case reflect.Float64:
		loadValue, errLoading := config.configureFileParsed.GetFloat64(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(loadValue))
	case reflect.Bool:
		loadValue, errLoading := config.configureFileParsed.GetBoolean(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(loadValue))
	case reflect.String:
		loadValue, errLoading := config.configureFileParsed.GetString(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(loadValue))
	default:
		return NewGoStructorNoValue(context.Value.Interface(), errors.New("can not parsed inserted type in GetBaseType of configuration by hocon"))
	}
}

func (config *HoconConfig) getSliceFromHocon(context *structContext) GoStructorValue {
	configLoad, structValue := config.typeSafeLoadConfigFile(context)
	if !configLoad {
		return *structValue
	}
	path := config.getElementName(context)
	fmt.Println("[HOCON]: level: debug. get path from hocon: ", path)
	valueIndirect := reflect.Indirect(context.Value)
	setupSlice := reflect.MakeSlice(valueIndirect.Type(), 1, 1)
	fmt.Println("[HOCON]: level: debug. type of first element at slice: ", setupSlice.Index(0).Kind())
	switch setupSlice.Index(0).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16:
		return NewGoStructorNoValue(context.Value.Interface(), errors.New("not supported yet"))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return NewGoStructorNoValue(context.Value.Interface(), errors.New("not supported yet"))
	case reflect.Int32:
		neededValues, errLoading := config.configureFileParsed.GetInt32List(path)
		if errLoading != nil {
			fmt.Println("loading error: ", errLoading)
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(neededValues))
	case reflect.Int64:
		neededValues, errLoading := config.configureFileParsed.GetInt64List(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(neededValues))
	case reflect.Float32:
		neededValues, errLoading := config.configureFileParsed.GetFloat32List(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(neededValues))
	case reflect.Float64:
		neededValues, errLoading := config.configureFileParsed.GetFloat64List(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(neededValues))
	case reflect.Bool:
		neededValues, errLoading := config.configureFileParsed.GetBooleanList(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(neededValues))
	case reflect.String:
		needValues, errLoading := config.configureFileParsed.GetStringList(path)
		if errLoading != nil {
			return NewGoStructorNoValue(context.Value.Interface(), errLoading)
		}
		return NewGoStructorTrueValue(reflect.ValueOf(needValues))
	case reflect.Complex64, reflect.Complex128:
		return NewGoStructorNoValue(context.Value.Interface(), errors.New("not supported yet a complex values"))
	default:
		return NewGoStructorNoValue(context.Value.Interface(), errors.New("can not recognize inserted type"))
	}
}

func (config *HoconConfig) getMapFromHocon(context *structContext) GoStructorValue {
	configLoad, structValue := config.typeSafeLoadConfigFile(context)
	if !configLoad {
		return *structValue
	}
	valueIndirect := reflect.Indirect(context.Value)
	path := config.getElementName(context)
	fmt.Println("[HOCON]: level: debuf. current type: ", valueIndirect.Kind().String())
	fmt.Println("[HOCON]: level: debuf.key of map: ", valueIndirect.Type().Key().Kind())
	fmt.Println("[HOCON]: level: debuf.value of map: ", valueIndirect.Type().Elem().Kind())
	getValue, errLoading := config.configureFileParsed.GetValue(path)
	if errLoading != nil {
		return NewGoStructorNoValue(context.Value.Interface(), errLoading)
	}
	object, errLoad := getValue.GetObject()
	if errLoad != nil {
		return NewGoStructorNoValue(context.Value.Interface(), errLoad)
	}
	keys := object.GetKeys()
	makebleMap := reflect.MakeMapWithSize(valueIndirect.Type(), 0)
	for _, key := range keys {
		value := config.GetBaseType(&structContext{
			StructField: reflect.StructField{
				Name: key,
			},
			Prefix: context.Prefix + "." + key,
			Value:  reflect.New(valueIndirect.Type().Elem()),
		})

		parsedKey := parseMapType(valueIndirect.Type().Key(), reflect.ValueOf(key))
		if parsedKey.CheckIsValue() {
			if value.CheckIsValue() {
				makebleMap.SetMapIndex(reflect.Indirect(parsedKey.Value).Convert(valueIndirect.Type().Key()), value.Value)
			} else {
				return NewGoStructorNoValue(parsedKey.Value, errors.New("can not parsed value for map"))
			}
		} else {
			return NewGoStructorNoValue(parsedKey.Value, errors.New("can not parsed key for map"))
		}

	}
	return NewGoStructorTrueValue(makebleMap)
}

func parseMapType(needType reflect.Type, value reflect.Value) GoStructorValue {
	switch needType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsedValue, errParsed := strconv.ParseInt(reflect.Indirect(value).Interface().(string), 10, 64)
		if errParsed != nil {
			return NewGoStructorNoValue(nil, errParsed)
		}

		return NewGoStructorTrueValue(reflect.ValueOf(parsedValue))
	case reflect.String:
		return NewGoStructorTrueValue(value)
	case reflect.Float32:
		parsedValue, errParsed := strconv.ParseFloat(reflect.Indirect(value).Interface().(string), 32)
		if errParsed != nil {
			return NewGoStructorNoValue(nil, errParsed)
		}

		return NewGoStructorTrueValue(reflect.ValueOf(parsedValue))
	case reflect.Float64:
		parsedValue, errParsed := strconv.ParseFloat(reflect.Indirect(value).Interface().(string), 64)
		if errParsed != nil {
			return NewGoStructorNoValue(nil, errParsed)
		}

		return NewGoStructorTrueValue(reflect.ValueOf(parsedValue))
	default:
		return NewGoStructorNoValue(nil, errors.New("can not set for map key by insert type: "+needType.Kind().String()))
	}
}
