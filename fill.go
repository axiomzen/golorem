package lorem

import (
	"errors"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

var ErrInvalidSpecification = errors.New("must provide a struct pointer")

// A ParseError occurs when an environment variable cannot be converted to
// the type required by a struct field during assignment.
type ParseError struct {
	Message   string
	FieldName string
	TypeName  string
	Tag       string
}

// A Decoder is a type that knows how to de-serialize environment variables
// into itself.
type LoremDecoder interface {
	LoremDecode(tag string) error
}

func Fill(spec interface{}) error {
	// must be a struct pointer
	s := reflect.ValueOf(spec)
	if s.Kind() != reflect.Ptr {
		return ErrInvalidSpecification
	}
	s = s.Elem()
	if s.Kind() != reflect.Struct {
		return ErrInvalidSpecification
	}
	typeOfSpec := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		// check for our tags
		loremTag := typeOfSpec.Field(i).Tag.Get("lorem")
		if !f.CanSet() || loremTag == "-" {
			// ignore this field
			continue
		}
		// check for embedded anonymous structs
		// ignored anonymous structs already covered, see here: https://play.golang.org/p/2FWYoLzWCV
		if typeOfSpec.Field(i).Anonymous && f.Kind() == reflect.Struct {
			embeddedPtr := f.Addr().Interface()
			if err := Fill(embeddedPtr); err != nil {
				return err
			}
			// populate the field itself
			f.Set(reflect.ValueOf(embeddedPtr).Elem())
		}

		err := processField(loremTag, f)
		if err != nil {
			return &ParseError{
				Message:   err.Error(),
				FieldName: typeOfSpec.Field(i).Name,
				TypeName:  f.Type().String(),
				Tag:       loremTag,
			}
		}
	}
	return nil
}

func processField(tag string, field reflect.Value) error {
	typ := field.Type()

	decoder := decoderFrom(field)
	if decoder != nil {
		return decoder.LoremDecoder(tag)
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
				} else {
					return errors.New("must have another thing after comma")
				}
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
				field.SetString(Word(min, max))
			case "sentence":
				field.SetString(Sentence(min, max))
			case "paragraph":
				field.SetString(Paragraph(min, max))
			case "url":
				field.SetString(Url())
			case "readableurl":
				field.SetString(ReadableUrl(Sentence(min, max)))
			case "host":
				field.SetString(Host())
			case "email":
				field.SetString(Email())
			}
		}
	case reflect.Int, reflect.Int64:
		field.SetInt(rand.Int())
	case reflect.Int32:
		field.SetInt(rand.Int31())
	case reflect.Int8:
		field.SetInt(IntRange(0, math.MaxInt8))
	case reflect.Int16:
		field.SetInt(IntRange(0, math.MaxInt16))
	case reflect.Uint, reflect.Uint64, reflect.Uint32:
		field.SetUint(rand.Uint32())
	case reflect.Uint8:
		field.SetUint(IntRange(0, math.MaxUInt8))
	case reflect.Uint16:
		field.SetUint(IntRange(0, math.MaxUInt16))
	case reflect.Bool:
		field.SetBool(rand.Int()%2 == 0)
	case reflect.Float32:
		field.SetFloat(rand.Float32())
	case reflect.Float64:
		field.SetFloat(rand.Float64())
	case reflect.Slice:
		//vals := strings.Split(value, ",")
		// make a random slice length?
		size := IntRange(0, 10)
		sl := reflect.MakeSlice(typ, size, size)
		for i, _ := range sl {
			err := processField(tag, sl.Index(i))
			if err != nil {
				return err
			}
		}
		field.Set(sl)
	}

	return nil
}

func decoderFrom(field reflect.Value) LoremDecoder {
	if field.CanInterface() {
		dec, ok := field.Interface().(LoremDecoder)
		if ok {
			return dec
		}
	}

	// also check if pointer-to-type implements Decoder,
	// and we can get a pointer to our field
	if field.CanAddr() {
		field = field.Addr()
		dec, ok := field.Interface().(LoremDecoder)
		if ok {
			return dec
		}
	}

	return nil
}
