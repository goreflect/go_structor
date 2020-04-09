# Gostructor [![Actions Status](https://github.com/goreflect/gostructor/workflows/CI_dev/badge.svg)](https://github.com/goreflect/gostructor/actions?query=workflow%3ACI_dev) [![Go Report Card](https://goreportcard.com/badge/github.com/goreflect/gostructor)](https://goreportcard.com/report/github.com/goreflect/gostructor) [![codecov](https://codecov.io/gh/goreflect/gostructor/branch/master/graph/badge.svg)](https://codecov.io/gh/goreflect/gostructor)
____
### Version: 0.4.1

Universal configuration library by tags

## Current supporting input formats

- hocon values
- default values
- environment variables

## Current supporting types

- int, int8, int16, int32, int64
- float32, float64
- string
- bool
- map[string\int]string\int\float32\float64
- slices of any types from (int32, int64, int, string, bool, float32, float64)

### Tags

- [x] cf_hocon - setup value for this field from hocon
- [x] cf_default - setup default value for this field
- [x] cf_env - setup value from env variable by name in this tag
- [ ] cf_yaml - setup value for this field from yaml (version > 0.6)
- [ ] cf_json - setup value for this field from json (version > 0.5)
- [ ] cf_server - setup value from configuration server like spring cloud config server or others (version>0.7)
- [ ] cf_vault - setup secret for this field from hashi corp vault (version>0.8)

## Running configuring by smart variant

For Run configuration by smart variant autostart analysing of using tags. library will start configuring  your structure by pipeline with all founded tags.

```go
type Test struct {
    MyValue1 string `cf_default:"turur" cf_hocon:"mySourceValue1"`
    MySlice1 []bool `cf_default:"true,false,false,true" cf_env:"MY_SIGNALS"`
}

// in this example do use 3 tags: cf_default (using default values which setup inline tag)
// cf_env - using environment variable
// cf_hocon - using hocon source file 

//....

func myConfigurator() {
    myStruct, errConfiguring := gostructir.ConfigureSmart(&Test{}, "testhocon.hocon")
    // check errConfiguring for any errors
    if errConfiguring != nil {
        /// action for error
    }

    // cast interface{} into Test structure
    myValues := myStruct.(*Test)
    // now, u structure already filled
} 

```

## Running configuring by setup

You can also setting configuring pipeline like this:

```go
type Test struct {
    MyValue1 string `cf_default:"turur" cf_hocon:"mySourceValue1"`
    MySlice1 []bool `cf_default:"true,false,false,true" cf_env:"MY_SIGNALS"`
}

func myConfigurator() {
    myStruct, errConfiguring := gostructir.ConfigureSetup(&Test{}, "testhocon.hocon", []infra.FuncType{
        infra.FunctionSetupEnvironment,
    })// you should setup only by order configure
    // check errConfiguring for any errors
    if errConfiguring != nil {
        /// action for error
    }

    // cast interface{} into Test structure
    myValues := myStruct.(*Test)
    // now, u structure already filled
} 

```
