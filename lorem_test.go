// Copyright 2012 Derek A. Rhodes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lorem

import "testing"
import "log"

type SimpleStruct struct {
	Int8               int8
	Int16              int16
	Int32              int32
	Int64              int64
	Int                int
	UInt8              uint8
	UInt16             uint16
	UInt32             uint32
	UInt64             uint64
	UInt               uint
	PlainString        string
	Word               string `lorem:"word"`
	WordWithRange      string `lorem:"word,10,11"`
	Sentence           string `lorem:"sentence"`
	SentenceWithRange  string `lorem:"sentence,10,11"`
	Paragraph          string `lorem:"paragraph"`
	ParagraphWithRange string `lorem:"paragraph,10,11"`
	URL                string `lorem:"url"`
	ReadableUrl        string `lorem:"readableurl"`
	Host               string `lorem:"host"`
	Email              string `lorem:"email"`
	Bool               bool
}

type StructWithPointers struct {
	Int8Pointer               *int8
	Int16Pointer              *int16
	Int32Pointer              *int32
	Int64Pointer              *int64
	IntPointer                *int
	UInt8Pointer              *uint8
	UInt16Pointer             *uint16
	UInt32Pointer             *uint32
	UInt64Pointer             *uint64
	UIntPointer               *uint
	PlainStringPointer        *string
	WordPointer               *string `lorem:"word"`
	WordWithRangePointer      *string `lorem:"word,10,11"`
	SentencePointer           *string `lorem:"sentence"`
	SentenceWithRangePointer  *string `lorem:"sentence,10,11"`
	ParagraphPointer          *string `lorem:"paragraph"`
	ParagraphWithRangePointer *string `lorem:"paragraph,10,11"`
	URLPointer                *string `lorem:"url"`
	ReadableUrlPointer        *string `lorem:"readableurl"`
	HostPointer               *string `lorem:"host"`
	EmailPointer              *string `lorem:"email"`
	BoolPointer               *bool
}

type OtherStruct struct {
	SubEmailPointer  *string `lorem:"email"`
	SubWordWithRange string  `lorem:"word,10,11"`
}

type StructWithSlices struct {
	// try a slice of things
	// tag applies to each thing, rather than the slice as a whole
	Words []string `lorem:"word"`
}

type StructWithSlicesOfPointers struct {
	// try a slice of pointers to things
	Sentences []*string `lorem:"sentence,10,11"`
}

// try a map
type StructWithMap struct {

	// should ignore it?
	Map        map[string]string
	MapWithTag map[string]string `lorem:"word"`
}

type StructWithStruct struct {
	OtherStruct        OtherStruct
	OtherStructPointer *OtherStruct
}

type StructWithEmbeddedStruct struct {
	OtherStruct
	URLPointer *string `lorem:"url"`
}

type StructWithEmbeddedStructPointer struct {
	*OtherStruct
	Word string `lorem:"word"`
}

type StructWithSliceOfStructs struct {
	OtherStructs        []OtherStruct
	OtherStructPointers []*OtherStruct
}

type StructWithIgnoredFields struct {
	IgnoredInt           int          `lorem:"-"`
	IgnoredUInt          uint         `lorem:"-"`
	IgnoredString        string       `lorem:"-"`
	IgnoredBool          bool         `lorem:"-"`
	IgnoredIntPointer    *int         `lorem:"-"`
	IgnoredUIntPointer   *uint        `lorem:"-"`
	IgnoredStringPointer *string      `lorem:"-"`
	IgnoredBoolPointer   *bool        `lorem:"-"`
	IgnoredStruct        OtherStruct  `lorem:"-"`
	IgnoredStructPointer *OtherStruct `lorem:"-"`
}

type StructWithIgnoredEmbeddedStruct struct {
	OtherStruct `lorem:"-"`
	Word        string `lorem:"word,10,11"`
}

func TestAll(t *testing.T) {
	for i := 1; i < 14; i++ {
		log.Print(word(i))
		for j := 1; j < 14; j++ {
			log.Print(Word(i, j))
			log.Print(Sentence(i, j))
			log.Print(Paragraph(i, j))
		}
		log.Print(Url())
		log.Print(Host())
		log.Print(Email())
	}
}
