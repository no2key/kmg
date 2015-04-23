package kmgYaml

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/bronze1man/kmg/kmgTest"
)

var unmarshalIntTest = 123
var unmarshalTimeTest = time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
var unmarshalTests = []struct {
	data  string
	value interface{}
}{
	{
		"",
		&struct{}{},
	}, {
		"{}", &struct{}{},
	}, {
		"v: hi",
		map[string]string{"v": "hi"},
	}, {
		"v: hi", map[string]interface{}{"v": "hi"},
	}, {
		"v: true",
		map[string]string{"v": "true"},
	}, {
		"v: true",
		map[string]interface{}{"v": true},
	}, {
		"v: 10",
		map[string]interface{}{"v": 10},
	}, {
		"v: 0b10",
		map[string]interface{}{"v": 2},
	}, {
		"v: 0xA",
		map[string]interface{}{"v": 10},
	}, {
		"v: 4294967296",
		map[string]int64{"v": 4294967296},
	}, {
		"v: 0.1",
		map[string]interface{}{"v": 0.1},
	}, {
		"v: .1",
		map[string]interface{}{"v": 0.1},
	}, {
		"v: .Inf",
		map[string]interface{}{"v": math.Inf(+1)},
	}, {
		"v: -.Inf",
		map[string]interface{}{"v": math.Inf(-1)},
	}, {
		"v: -10",
		map[string]interface{}{"v": -10},
	}, {
		"v: -.1",
		map[string]interface{}{"v": -0.1},
	},

	// Simple values.
	{
		"123",
		&unmarshalIntTest,
	},

	// Floats from spec
	{
		"canonical: 6.8523e+5",
		map[string]interface{}{"canonical": 6.8523e+5},
	}, {
		"expo: 685.230_15e+03",
		map[string]interface{}{"expo": 685.23015e+03},
	}, {
		"fixed: 685_230.15",
		map[string]interface{}{"fixed": 685230.15},
	}, {
		"neginf: -.inf",
		map[string]interface{}{"neginf": math.Inf(-1)},
	}, {
		"fixed: 685_230.15",
		map[string]float64{"fixed": 685230.15},
	},
	//{"sexa: 190:20:30.15", map[string]interface{}{"sexa": 0}}, // Unsupported
	//{"notanum: .NaN", map[string]interface{}{"notanum": math.NaN()}}, // Equality of NaN fails.

	// Bools from spec
	{
		"canonical: y",
		map[string]interface{}{"canonical": true},
	}, {
		"answer: NO",
		map[string]interface{}{"answer": false},
	}, {
		"logical: True",
		map[string]interface{}{"logical": true},
	}, {
		"option: on",
		map[string]interface{}{"option": true},
	}, {
		"option: on",
		map[string]bool{"option": true},
	},
	// Ints from spec
	{
		"canonical: 685230",
		map[string]interface{}{"canonical": 685230},
	}, {
		"decimal: +685_230",
		map[string]interface{}{"decimal": 685230},
	}, {
		"octal: 02472256",
		map[string]interface{}{"octal": 685230},
	}, {
		"hexa: 0x_0A_74_AE",
		map[string]interface{}{"hexa": 685230},
	}, {
		"bin: 0b1010_0111_0100_1010_1110",
		map[string]interface{}{"bin": 685230},
	}, {
		"bin: -0b101010",
		map[string]interface{}{"bin": -42},
	}, {
		"decimal: +685_230",
		map[string]int{"decimal": 685230},
	},

	//{"sexa: 190:20:30", map[string]interface{}{"sexa": 0}}, // Unsupported

	// Nulls from spec
	{
		"empty:",
		map[string]interface{}{"empty": nil},
	}, {
		"canonical: ~",
		map[string]interface{}{"canonical": nil},
	}, {
		"english: null",
		map[string]interface{}{"english": nil},
	}, {
		"~: null key",
		map[interface{}]string{nil: "null key"},
	}, {
		"empty:",
		map[string]*bool{"empty": nil},
	},

	// Flow sequence
	{
		"seq: [A,B]",
		map[string]interface{}{"seq": []interface{}{"A", "B"}},
	}, {
		"seq: [A,B,C,]",
		map[string][]string{"seq": {"A", "B", "C"}},
	}, {
		"seq: [A,1,C]",
		map[string][]string{"seq": {"A", "1", "C"}},
	}, {
		"seq: [A,1,C]",
		map[string][]int{"seq": {1}},
	}, {
		"seq: [A,1,C]",
		map[string]interface{}{"seq": []interface{}{"A", 1, "C"}},
	},
	// Block sequence
	{
		"seq:\n - A\n - B",
		map[string]interface{}{"seq": []interface{}{"A", "B"}},
	}, {
		"seq:\n - A\n - B\n - C",
		map[string][]string{"seq": {"A", "B", "C"}},
	}, {
		"seq:\n - A\n - 1\n - C",
		map[string][]string{"seq": {"A", "1", "C"}},
	}, {
		"seq:\n - A\n - 1\n - C",
		map[string][]int{"seq": {1}},
	}, {
		"seq:\n - A\n - 1\n - C",
		map[string]interface{}{"seq": []interface{}{"A", 1, "C"}},
	},

	// Literal block scalar
	{
		"scalar: | # Comment\n\n literal\n\n \ttext\n\n",
		map[string]string{"scalar": "\nliteral\n\n\ttext\n"},
	},

	// Folded block scalar
	{
		"scalar: > # Comment\n\n folded\n line\n \n next\n line\n  * one\n  * two\n\n last\n line\n\n",
		map[string]string{"scalar": "\nfolded line\nnext line\n * one\n * two\n\nlast line\n"},
	},

	// Map inside interface with no type hints.
	{
		"a: {b: c}",
		map[string]interface{}{"a": map[interface{}]interface{}{"b": "c"}},
	},

	// Structs and type conversions.
	{
		"Hello: world",
		&struct{ Hello string }{"world"},
	}, {
		"A: {B: c}",
		&struct{ A struct{ B string } }{struct{ B string }{"c"}},
	}, {
		"A: {B: c}",
		&struct{ A *struct{ B string } }{&struct{ B string }{"c"}},
	}, {
		"A: {b: c}",
		&struct{ A map[string]string }{map[string]string{"b": "c"}},
	}, {
		"A: {b: c}",
		&struct{ A *map[string]string }{&map[string]string{"b": "c"}},
	}, {
		"A:",
		&struct{ A map[string]string }{},
	}, {
		"A: 1",
		&struct{ A int }{1},
	}, {
		"A: [1, 2]",
		&struct{ A []int }{[]int{1, 2}},
	}, {
		"A: [1, 2]",
		&struct{ A [2]int }{[2]int{1, 2}},
	}, {
		"A: 1",
		&struct{ B int }{0},
	}, {
		"a: 1",
		&struct {
			B int "a"
		}{1},
	}, {
		"A: y",
		&struct{ A bool }{true},
	},

	// Some cross type conversions
	{
		"v: 42",
		map[string]uint{"v": 42},
	}, {
		"v: -42",
		map[string]uint{},
	}, {
		"v: 4294967296",
		map[string]uint64{"v": 4294967296},
	}, {
		"v: -4294967296",
		map[string]uint64{},
	},

	// Overflow cases.
	{
		"v: 4294967297",
		map[string]int32{},
	}, {
		"v: 128",
		map[string]int8{},
	},

	// Quoted values.
	{
		"'1': '\"2\"'",
		map[interface{}]interface{}{"1": "\"2\""},
	}, {
		"v:\n- A\n- 'B\n\n  C'\n",
		map[string][]string{"v": {"A", "B\nC"}},
	},

	// Explicit tags.
	{
		"v: !!float '1.1'",
		map[string]interface{}{"v": 1.1},
	}, {
		"v: !!null ''",
		map[string]interface{}{"v": nil},
	}, {
		"%TAG !y! tag:yaml.org,2002:\n---\nv: !y!int '1'",
		map[string]interface{}{"v": 1},
	},

	// Anchors and aliases.
	{
		"A: &x 1\nB: &y 2\nC: *x\nD: *y\n",
		&struct{ A, B, C, D int }{1, 2, 1, 2},
	}, {
		"A: &a {C: 1}\nB: *a",
		&struct {
			A, B struct {
				C int
			}
		}{struct{ C int }{1}, struct{ C int }{1}},
	}, {
		"A: &a [1, 2]\nB: *a",
		&struct{ B []int }{[]int{1, 2}},
	},

	// Bug #1133337
	{
		"foo: ''",
		map[string]*string{"foo": new(string)},
	}, {
		"foo: null",
		map[string]string{},
	},

	// Ignored field
	{
		"A: 1\nB: 2\n",
		&struct {
			A int
			B int "-"
		}{1, 0},
	},

	// Bug #1191981
	{
		"" +
			"%YAML 1.1\n" +
			"--- !!str\n" +
			`"Generic line break (no glyph)\n\` + "\n" +
			` Generic line break (glyphed)\n\` + "\n" +
			` Line separator\u2028\` + "\n" +
			` Paragraph separator\u2029"` + "\n",
		"" +
			"Generic line break (no glyph)\n" +
			"Generic line break (glyphed)\n" +
			"Line separator\u2028Paragraph separator\u2029",
	},

	// Struct inlining
	{
		"A: 1\nB: 2\nC: 3\n",
		&struct {
			A int
			C inlineB `yaml:",inline"`
		}{1, inlineB{2, inlineC{3}}},
	},
	{
		`0: 2
1: 0.4
2: 0.7
3: 1
4: 1.5
`,
		map[int]float64{
			0: 2.0,
			1: 0.4,
			2: 0.7,
			3: 1.0,
			4: 1.5,
		},
	},
	{
		`2001-02-03T04:05:06Z`,
		&unmarshalTimeTest,
	},
}

type inlineB struct {
	B       int
	inlineC `yaml:",inline"`
}

type inlineC struct {
	C int
}

func (c *S) TestUnmarshal() {
	for _, item := range unmarshalTests {
		t := reflect.ValueOf(item.value).Type()
		var value interface{}
		switch t.Kind() {
		case reflect.Map:
			value = reflect.MakeMap(t).Interface()
		case reflect.String:
			t := reflect.ValueOf(item.value).Type()
			v := reflect.New(t)
			value = v.Interface()
		case reflect.Ptr:
			pt := t
			pv := reflect.New(pt.Elem())
			value = pv.Interface()
		default:
			pt := t
			pv := reflect.New(pt)
			value = pv.Interface()
		}
		err := Unmarshal([]byte(item.data), value)
		c.Equal(err, nil)
		//c.Assert(err, IsNil, Commentf("Item #%d", i))
		if t.Kind() == reflect.String {
			c.Equal(*value.(*string), item.value)
			//c.Assert(*value.(*string), Equals, item.value, Commentf("Item #%d", i))
		} else {
			c.Equal(value, item.value)
			//c.Assert(value, DeepEquals, item.value, Commentf("Item #%d", i))
		}
	}
}

func TestUnmarshalNaN(ot *testing.T) {
	c := kmgTest.NewTestTools(ot)
	value := map[string]interface{}{}
	err := Unmarshal([]byte("notanum: .NaN"), &value)
	c.Equal(err, nil)
	c.Equal(math.IsNaN(value["notanum"].(float64)), true)
	//c.Assert(err, IsNil)
	//c.Assert(math.IsNaN(value["notanum"].(float64)), Equals, true)
}

var unmarshalErrorTests = []struct {
	data, error string
}{
	{"v: !!float 'error'", "YAML error: Can't decode !!str 'error' as a !!float"},
	{"v: [A,", "YAML error: line 1: did not find expected node content"},
	{"v:\n- [A,", "YAML error: line 2: did not find expected node content"},
	{"a: *b\n", "YAML error: Unknown anchor 'b' referenced"},
	{"a: &a\n  b: *a\n", "YAML error: Anchor 'a' value contains itself"},
}

func TestUnmarshalErrors(ot *testing.T) {
	c := kmgTest.NewTestTools(ot)
	for _, item := range unmarshalErrorTests {
		var value interface{}
		err := Unmarshal([]byte(item.data), &value)
		c.Equal(err.Error(), item.error)
	}
}

var setterTests = []struct {
	data, tag string
	value     interface{}
}{
	{"_: {hi: there}", "!!map", map[interface{}]interface{}{"hi": "there"}},
	{"_: [1,A]", "!!seq", []interface{}{1, "A"}},
	{"_: 10", "!!int", 10},
	{"_: null", "!!null", nil},
	{"_: !!foo 'BAR!'", "!!foo", "BAR!"},
}

var setterResult = map[int]bool{}

type typeWithSetter struct {
	tag   string
	value interface{}
}

func (o *typeWithSetter) SetYAML(tag string, value interface{}) (ok bool) {
	o.tag = tag
	o.value = value
	if i, ok := value.(int); ok {
		if result, ok := setterResult[i]; ok {
			return result
		}
	}
	return true
}

type typeWithSetterField struct {
	Field *typeWithSetter "_"
}

func TestUnmarshalWithSetter(ot *testing.T) {
	c := kmgTest.NewTestTools(ot)
	for _, item := range setterTests {
		obj := &typeWithSetterField{}
		err := Unmarshal([]byte(item.data), obj)
		c.Equal(err, nil)
		c.Ok(obj.Field != nil)
		c.Equal(obj.Field.tag, item.tag)
		c.Equal(obj.Field.value, item.value)
		/*
			c.Assert(err, IsNil)
			c.Assert(obj.Field, NotNil,
				Commentf("Pointer not initialized (%#v)", item.value))
			c.Assert(obj.Field.tag, Equals, item.tag)
			c.Assert(obj.Field.value, DeepEquals, item.value)
		*/
	}
}

func TestUnmarshalWholeDocumentWithSetter(ot *testing.T) {
	c := kmgTest.NewTestTools(ot)
	obj := &typeWithSetter{}
	err := Unmarshal([]byte(setterTests[0].data), obj)
	c.Equal(err, nil)
	c.Equal(obj.tag, setterTests[0].tag)
	//c.Assert(err, IsNil)
	//c.Assert(obj.tag, Equals, setterTests[0].tag)

	value, ok := obj.value.(map[interface{}]interface{})
	c.Equal(ok, true)
	c.Equal(value["_"], setterTests[0].value)

	//c.Assert(ok, Equals, true)
	//c.Assert(value["_"], DeepEquals, setterTests[0].value)
}

func TestUnmarshalWithFalseSetterIgnoresValue(ot *testing.T) {
	c := kmgTest.NewTestTools(ot)
	setterResult[2] = false
	setterResult[4] = false
	defer func() {
		delete(setterResult, 2)
		delete(setterResult, 4)
	}()

	m := map[string]*typeWithSetter{}
	data := "{abc: 1, def: 2, ghi: 3, jkl: 4}"
	err := Unmarshal([]byte(data), m)
	c.Equal(err, nil)
	c.Ok(m["abc"] != nil)
	c.Equal(m["def"], nil)
	c.Ok(m["ghi"] != nil)
	c.Equal(m["jkl"], nil)
	c.Equal(m["abc"].value, 1)
	c.Equal(m["ghi"].value, 3)
	/*
		c.Assert(err, IsNil)
		c.Assert(m["abc"], NotNil)
		c.Assert(m["def"], IsNil)
		c.Assert(m["ghi"], NotNil)
		c.Assert(m["jkl"], IsNil)

		c.Assert(m["abc"].value, Equals, 1)
		c.Assert(m["ghi"].value, Equals, 3)
	*/
}

func TestUnmarshalTypeNotMatch(t *testing.T) {
	data := `t1:
     k1: v1`
	out := map[string][]map[string]string{}
	err := Unmarshal([]byte(data), &out)
	kmgTest.Ok(err != nil)
}

//var data []byte
//func init() {
//	var err error
//	data, err = ioutil.ReadFile("/tmp/file.yaml")
//	if err != nil {
//		panic(err)
//	}
//}
//
//func (s *S) BenchmarkUnmarshal(c *C) {
//	var err error
//	for i := 0; i < c.N; i++ {
//		var v map[string]interface{}
//		err = goyaml.Unmarshal(data, &v)
//	}
//	if err != nil {
//		panic(err)
//	}
//}
//
//func (s *S) BenchmarkMarshal(c *C) {
//	var v map[string]interface{}
//	goyaml.Unmarshal(data, &v)
//	c.ResetTimer()
//	for i := 0; i < c.N; i++ {
//		goyaml.Marshal(&v)
//	}
//}
