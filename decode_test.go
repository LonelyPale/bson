package bson

import (
	"errors"
	"reflect"
	"testing"
)

type element struct {
	typ     byte
	ename   []byte
	element []byte
}

var decodeTests = []struct {
	bson     []byte
	expected []element
	err      error
}{{
	// test1.bson
	bson: []byte{0x0e, 0x00, 0x00, 0x00, 0x10, 0x69, 0x6e, 0x74, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
	expected: []element{{
		typ:     0x10,
		ename:   cstring("int"),
		element: []byte{0x01, 0, 0, 0},
	}},
}, {
	// test2.bson
	bson: []byte{0x14, 0x00, 0x00, 0x00, 0x12, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	expected: []element{{
		typ:     0x12,
		ename:   cstring("int64"),
		element: []byte{0x01, 0, 0, 0, 0, 0, 0, 0},
	}},
}, {
	// test3.bson
	bson: []byte{0x15, 0x00, 0x00, 0x00, 0x01, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x00, 0x2b, 0x87, 0x16, 0xd9, 0xce, 0xf7, 0xf1, 0x3f, 0x00},
	expected: []element{{
		typ:     0x1,
		ename:   cstring("double"),
		element: []uint8{0x2b, 0x87, 0x16, 0xd9, 0xce, 0xf7, 0xf1, 0x3f},
	}},
}, {
	// test4.bson
	bson: []byte{0x12, 0x00, 0x00, 0x00, 0x09, 0x75, 0x74, 0x63, 0x00, 0x0b, 0x98, 0x8c, 0x2b, 0x33, 0x01, 0x00, 0x00, 0x00},
	expected: []element{{
		typ:     0x09, // timestamp
		ename:   cstring("utc"),
		element: []uint8{0xb, 0x98, 0x8c, 0x2b, 0x33, 0x1, 0x0, 0x0},
	}},
}, {
	// test5.bson
	bson: []byte{0x1d, 0x00, 0x00, 0x00, 0x02, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x00, 0x00},
	expected: []element{{
		typ:     0x02, // utf-8 string
		ename:   cstring("string"),
		element: []byte{0x73, 0x6f, 0x6d, 0x65, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x00},
	}},
}, {
	// test6.bson
	bson: []byte{0x40, 0x00, 0x00, 0x00, 0x04, 0x61, 0x72, 0x72, 0x61, 0x79, 0x5b, 0x69, 0x6e, 0x74, 0x5d, 0x00, 0x2f, 0x00, 0x00, 0x00, 0x10, 0x30, 0x00, 0x01, 0x00, 0x00, 0x00, 0x10, 0x31, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10, 0x32, 0x00, 0x03, 0x00, 0x00, 0x00, 0x10, 0x33, 0x00, 0x04, 0x00, 0x00, 0x00, 0x10, 0x34, 0x00, 0x05, 0x00, 0x00, 0x00, 0x10, 0x35, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00},
	expected: []element{{
		typ:     0x04,
		ename:   cstring("array[int]"),
		element: []byte("/\x00\x00\x00\x100\x00\x01\x00\x00\x00\x101\x00\x02\x00\x00\x00\x102\x00\x03\x00\x00\x00\x103\x00\x04\x00\x00\x00\x104\x00\x05\x00\x00\x00\x105\x00\x06\x00\x00\x00\x00"),
	}},
}, {
	// test7.bson
	bson: []byte{0x2f, 0x00, 0x00, 0x00, 0x04, 0x61, 0x72, 0x72, 0x61, 0x79, 0x5b, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x5d, 0x00, 0x1b, 0x00, 0x00, 0x00, 0x01, 0x30, 0x00, 0x2b, 0x87, 0x16, 0xd9, 0xce, 0xf7, 0xf1, 0x3f, 0x01, 0x31, 0x00, 0x96, 0x43, 0x8b, 0x6c, 0xe7, 0xfb, 0x00, 0x40, 0x00, 0x00},
	expected: []element{{
		typ:     0x04,
		ename:   cstring("array[double]"),
		element: []byte("\x1b\x00\x00\x00\x010\x00+\x87\x16\xd9\xce\xf7\xf1?\x011\x00\x96C\x8bl\xe7\xfb\x00@\x00"),
	}},
}, {
	// test8.bson
	bson: []byte{0x1d, 0x00, 0x00, 0x00, 0x03, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x00, 0x0e, 0x00, 0x00, 0x00, 0x10, 0x69, 0x6e, 0x74, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00},
	expected: []element{{
		typ:     0x03,
		ename:   cstring("document"),
		element: []byte("\x0e\x00\x00\x00\x10int\x00\x01\x00\x00\x00\x00"),
	}},
}, {
	// test9.bson
	bson: []byte{0x0b, 0x00, 0x00, 0x00, 0x0a, 0x6e, 0x75, 0x6c, 0x6c, 0x00, 0x00},
	expected: []element{{
		typ:     0x0a,
		ename:   cstring("null"),
		element: nil,
	}},
}, {
	// test10.bson
	bson: []byte{0x13, 0x00, 0x00, 0x00, 0x0b, 0x72, 0x65, 0x67, 0x65, 0x78, 0x00, 0x31, 0x32, 0x33, 0x34, 0x00, 0x69, 0x00, 0x00},
	expected: []element{{
		typ:     0x0b,
		ename:   cstring("regex"),
		element: []byte("1234\x00i\x00"),
	}},
}, {
	// test11.bson
	bson: []byte("\x16\x00\x00\x00\x02hello\x00\x06\x00\x00\x00world\x00\x00"),
	expected: []element{{
		typ:     0x02, // utf-8 string
		ename:   cstring("hello"),
		element: []byte("world\x00"),
	}},
	err: nil,
}, {
	// test12.bson
	bson: []byte("\x31\x00\x00\x00\x04BSON\x00\x26\x00\x00\x00\x020\x00\x08\x00\x00\x00awesome\x00\x011\x00\x33\x33\x33\x33\x33\x33\x14\x40\x102\x00\xc2\x07\x00\x00\x00\x00"),
	expected: []element{{
		typ:     0x4, // bson array
		ename:   cstring("BSON"),
		element: []byte("\x26\x00\x00\x00\x020\x00\x08\x00\x00\x00awesome\x00\x011\x00\x33\x33\x33\x33\x33\x33\x14\x40\x102\x00\xc2\x07\x00\x00\x00"),
	}},
	err: nil,
}, {
	// test13.bson
	bson: []byte{0x23, 0x00, 0x00, 0x00, 0x04, 0x61, 0x72, 0x72, 0x61, 0x79, 0x5b, 0x62, 0x6f, 0x6f, 0x6c, 0x5d, 0x00, 0x11, 0x00, 0x00, 0x00, 0x08, 0x30, 0x00, 0x01, 0x08, 0x31, 0x00, 0x00, 0x08, 0x32, 0x00, 0x01, 0x00, 0x00},
	expected: []element{{
		typ:     0x04,
		ename:   cstring("array[bool]"),
		element: []byte("\x11\x00\x00\x00\b0\x00\x01\b1\x00\x00\b2\x00\x01\x00"),
	}},
}, {
	// test14.bson
	bson: []byte{0x33, 0x00, 0x00, 0x00, 0x04, 0x61, 0x72, 0x72, 0x61, 0x79, 0x5b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5d, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x02, 0x30, 0x00, 0x06, 0x00, 0x00, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x02, 0x31, 0x00, 0x06, 0x00, 0x00, 0x00, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x00, 0x00, 0x00},
	expected: []element{{
		typ:     0x04,
		ename:   cstring("array[string]"),
		element: []byte("\x1f\x00\x00\x00\x020\x00\x06\x00\x00\x00hello\x00\x021\x00\x06\x00\x00\x00world\x00\x00"),
	}},
}}

func TestReader(t *testing.T) {
	for _, tt := range decodeTests {
		r := reader{bson: tt.bson[4 : len(tt.bson)-1]}
		got := make([]element, 0)
		for r.Next() {
			typ, ename, value := r.Element()
			got = append(got, element{typ, ename, value})
		}
		err := r.Err()
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("bsonIter %v: expected err %v, got %v", tt.bson, tt.err, err)
			continue
		}
		if !reflect.DeepEqual(tt.expected, got) {
			t.Errorf("bsonIter %q: expected %#q, got %#q", tt.bson, tt.expected, got)
		}
	}
}

func TestDecodeMap(t *testing.T) {
	return
	for _, tt := range decodeTests {
		got := make(map[string]interface{})
		err := decode(tt.bson, &got)
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("decode(%q): expected err %v, got %v", tt.bson, tt.err, err)
			continue
		}
		if !reflect.DeepEqual(tt.expected, got) {
			t.Errorf("decode(%q): expected %q, got %q", tt.expected, got)
		}
	}
}

var readInt32Tests = []struct {
	data     []byte
	expected int
	rest     []byte
}{{
	data:     []byte{0x1, 0x1, 0x0, 0x0},
	expected: 0x101,
	rest:     []byte{},
}, {
	data:     []byte{0x0, 0x0, 0x0, 0x1},
	expected: 0x01000000,
	rest:     []byte{},
}, {
	data:     []byte{0x0f, 0x0f, 0x0f, 0x0f, 0x0f},
	expected: 0x0f0f0f0f,
	rest:     []byte{0x0f},
}}

func TestReadInt32(t *testing.T) {
	for _, tt := range readInt32Tests {
		got, rest := readInt32(tt.data)
		if got != tt.expected || !reflect.DeepEqual(tt.rest, rest) {
			t.Errorf("readInt32(%v): expected %v %v, got %v, %v", tt.data, tt.expected, tt.rest, got, rest)
		}
	}
}

func cstring(s string) []byte {
	return append([]byte(s), 0)
}

var readCstringTests = []struct {
	data           []byte
	expected, rest []byte
	err            error
}{{
	data:     []byte{},
	expected: nil,
	rest:     nil,
	err:      errors.New("bson: cstring missing \\0"),
}, {
	data:     cstring("bson"),
	expected: cstring("bson"),
	rest:     []byte{},
	err:      nil,
}, {
	data:     cstring("bson\x00"),
	expected: cstring("bson"),
	rest:     []byte{0},
	err:      nil,
}}

func TestReadCstring(t *testing.T) {
	for _, tt := range readCstringTests {
		got, rest, err := readCstring(tt.data)
		if !reflect.DeepEqual(tt.err, err) || !reflect.DeepEqual(tt.expected, got) || !reflect.DeepEqual(tt.rest, rest) {
			t.Errorf("readCstring(%v): expected %v %v %v, got %v %v %v", tt.data, tt.expected, tt.rest, tt.err, got, rest, err)
		}
	}
}