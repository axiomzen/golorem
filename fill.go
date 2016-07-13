package lorem

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

//
var errInvalidSpecification = errors.New("must provide a struct pointer")

// A ParseError occurs when an environment variable cannot be converted to
// the type required by a struct field during assignment.
type ParseError struct {
	Message   string
	FieldName string
	TypeName  string
	Tag       string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("envconfig.Process: error %s for fieldname %s: has type %s and tag %s", e.Message, e.FieldName, e.TypeName, e.Tag)
}

// A Decoder is a type that knows how to de-serialize environment variables
// into itself.
type Decoder interface {
	LoremDecode(tag string) error
}

// this will handle everything
func loremizeRec(loremTag string, field reflect.Value) error {

	if !field.CanSet() || loremTag == "-" {
		// ignore this field
		return nil
	}
	// check for decoder
	decoder := decoderFrom(field)
	if decoder != nil {
		return decoder.LoremDecode(loremTag)
	}

	// check for pointer first
	typ := field.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if field.IsNil() {
			field.Set(reflect.New(typ))
		}
		field = field.Elem()
	}

	switch field.Kind() {
	case reflect.Struct:
		// call loremizeRec on each field
		typeOfField := field.Type()
		//if
		//todo: field.Anonymous
		for i := 0; i < field.NumField(); i++ {
			subField := field.Field(i)
			err := loremizeRec(typeOfField.Field(i).Tag.Get("lorem"), subField)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		// init slice, call loremizeRec on each slice entry
		size := IntRange(1, 10)
		sl := reflect.MakeSlice(typ, size, size)
		for i := 0; i < size; i++ {
			//err := processField(tag, sl.Index(i))
			sliceIndex := sl.Index(i)
			err := loremizeRec(loremTag, sliceIndex)
			if err != nil {
				return err
			}
		}
		field.Set(sl)
	default:
		// handle simple type
		err := processField(loremTag, field)
		if err != nil {
			return err
		}
	}
	return nil
}

// Loremize will fill in the structure with random stuff
// using lorme ipsum for strings
func Loremize(spec interface{}) error {
	// must be a struct pointer
	value := reflect.ValueOf(spec)
	if value.Kind() != reflect.Ptr {
		return errInvalidSpecification
	}
	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return errInvalidSpecification
	}
	typeOfValue := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		loremTag := typeOfValue.Field(i).Tag.Get("lorem")
		err := loremizeRec(loremTag, field)
		if err != nil {
			return &ParseError{
				Message:   err.Error(),
				FieldName: typeOfValue.Field(i).Name,
				TypeName:  field.Type().String(),
				Tag:       loremTag,
			}
		}

		// // check for our tags
		// loremTag := typeOfValue.Field(i).Tag.Get("lorem")
		// if !field.CanSet() || loremTag == "-" {
		// 	// ignore this field
		// 	continue
		// }

		// typeOfField := field.Type()
		// if field.Kind() == reflect.Ptr {
		// 	typeOfField = typeOfField.Elem()
		// 	if field.IsNil() {
		// 		field.Set(reflect.New(typeOfField))
		// 	}
		// 	field = field.Elem()
		// }
		// // check for structs
		// // ignored anonymous structs already covered, see here: https://play.golang.org/p/2FWYoLzWCV
		// if field.Kind() == reflect.Struct {
		// 	//if typeOfSpec.Field(i).Anonymous {
		// 	embeddedPtr := field.Addr().Interface()
		// 	if err := Loremize(embeddedPtr); err != nil {
		// 		return err
		// 	}
		// 	// populate the ptr field itself
		// 	field.Set(reflect.ValueOf(embeddedPtr).Elem())
		// 	//} else if err := Loremize(f.Addr()); err != nil {
		// 	//	return err
		// 	//}
		// 	continue
		// } else if field.Kind() == reflect.Slice {
		// 	// create the slice
		// 	// iterate over fields
		// 	// check for pointers

		// 	size := IntRange(1, 10)
		// 	sl := reflect.MakeSlice(typeOfField, size, size)
		// 	for i := 0; i < size; i++ {
		// 		//err := processField(tag, sl.Index(i))
		// 		sliceIndex := sl.Index(i)
		// 		// need function that takes anything (including structs)
		// 		typeOfSliceIndex := sliceIndex.Type()
		// 		if sliceIndex.Kind() == reflect.Ptr {
		// 			typeOfSliceIndex = typeOfSliceIndex.Elem()
		// 			if sliceIndex.IsNil() {
		// 				sliceIndex.Set(reflect.New(typeOfSliceIndex))
		// 			}
		// 			sliceIndex = sliceIndex.Elem()
		// 		}

		// 	}
		// 	field.Set(sl)
		// } else {
		// 	// handle simple field
		// }

		// /*
		// 			  // if a pointer to a struct is passed, get the type of the dereferenced object
		// 	  if typ.Kind() == reflect.Ptr{
		// 	    typ = typ.Elem()
		// 	  }
		// */

		// err := processField(loremTag, f)
		// if err != nil {
		// 	return &ParseError{
		// 		Message:   err.Error(),
		// 		FieldName: typeOfSpec.Field(i).Name,
		// 		TypeName:  f.Type().String(),
		// 		Tag:       loremTag,
		// 	}
		// }
	}
	return nil
}

func processField(tag string, field reflect.Value) error {
	typ := field.Type()

	decoder := decoderFrom(field)
	if decoder != nil {
		return decoder.LoremDecode(tag)
	}

	// handle pointers
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if field.IsNil() {
			field.Set(reflect.New(typ))
		}
		field = field.Elem()
	}

	// no lorem tag specified, use default for everything
	switch typ.Kind() {
	case reflect.String:
		if tag == "" {
			field.SetString(Word(2, 10))
		} else {
			args := strings.Split(tag, ",")
			if args[0] == "" {
				// just fill in nextone
				if len(args) > 1 {
					field.SetString(args[1])
					return nil
				}
				return errors.New("must have another thing after comma")
			}

			var min = int64(2)
			var max = int64(10)
			if len(args) == 3 {
				var err error
				min, err = strconv.ParseInt(args[1], 10, 32)
				if err != nil {
					return err
				}

				max, err = strconv.ParseInt(args[2], 10, 32)
				if err != nil {
					return err
				}
			}

			switch args[0] {
			case "word":
				field.SetString(Word(int(min), int(max)))
			case "sentence":
				field.SetString(Sentence(int(min), int(max)))
			case "paragraph":
				field.SetString(Paragraph(int(min), int(max)))
			case "url":
				field.SetString(URL())
			case "readablepath":
				field.SetString(ReadablePath(Sentence(int(min), int(max))))
			case "host":
				field.SetString(Host())
			case "email":
				field.SetString(Email())
			}
		}
	case reflect.Int, reflect.Int64:
		field.SetInt(int64(rand.Int()))
	case reflect.Int32:
		field.SetInt(int64(rand.Int31()))
	case reflect.Int8:
		field.SetInt(int64(IntRange(0, math.MaxInt8)))
	case reflect.Int16:
		field.SetInt(int64(IntRange(0, math.MaxInt16)))
	case reflect.Uint, reflect.Uint64, reflect.Uint32:
		field.SetUint(uint64(rand.Uint32()))
	case reflect.Uint8:
		field.SetUint(uint64(IntRange(0, math.MaxUint8)))
	case reflect.Uint16:
		field.SetUint(uint64(IntRange(0, math.MaxUint16)))
	case reflect.Bool:
		field.SetBool(rand.Int()%2 == 0)
	case reflect.Float32:
		field.SetFloat(float64(rand.Float32()))
	case reflect.Float64:
		field.SetFloat(rand.Float64())
		// case reflect.Slice:
		// 	// make a random slice length?
		// 	size := IntRange(1, 10)
		// 	sl := reflect.MakeSlice(typ, size, size)
		// 	for i := 0; i < size; i++ {
		// 		err := processField(tag, sl.Index(i))
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// 	field.Set(sl)
	default:
	}

	return nil
}

func decoderFrom(field reflect.Value) Decoder {
	if field.CanInterface() {
		dec, ok := field.Interface().(Decoder)
		if ok {
			return dec
		}
	}

	// also check if pointer-to-type implements Decoder,
	// and we can get a pointer to our field
	if field.CanAddr() {
		field = field.Addr()
		dec, ok := field.Interface().(Decoder)
		if ok {
			return dec
		}
	}

	return nil
}
