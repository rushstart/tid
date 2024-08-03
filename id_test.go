package tid

import (
	"io"
	"math/rand"
	rand2 "math/rand/v2"
	"reflect"
	"testing"
)

func testData() map[ID]string {
	data := make(map[ID]string, 8)
	data[From[rand.Rand]()] = From[rand.Rand]().String()
	data[From[rand2.Rand]()] = From[rand2.Rand]().String()
	data[From[*rand.Rand]()] = From[*rand.Rand]().String()
	data[From[io.Reader]()] = From[io.Reader]().String()
	data[From[*io.Reader]()] = From[*io.Reader]().String() // ?
	data[From[int]()] = From[int]().String()
	data[From[int]("some-tag")] = From[int]("some-tag").String()
	data[From[struct {
		r  rand.Rand  `tag:"r"`
		r2 rand2.Rand `tag:"r_2"`
	}]("some-tag")] = From[struct {
		r  rand.Rand  `tag:"r"`
		r2 rand2.Rand `tag:"r_2"`
	}]("some-tag").String()
	return data
}

func assertEqual(t *testing.T, expected, actual any) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Not equal:\nexpected: %v\nactual: %v", expected, actual)
	}
}

func assertTrue(t *testing.T, expected bool) {
	if !expected {
		t.Error("should be true")
	}
}

func TestFrom(t *testing.T) {
	data := testData()

	assertEqual(t, data[From[rand.Rand]()], "math/rand/Rand")
	assertEqual(t, data[From[rand2.Rand]()], "math/rand/v2/Rand")
	assertEqual(t, data[From[*rand.Rand]()], "*math/rand/Rand")
	assertEqual(t, data[From[io.Reader]()], "io/Reader")
	assertEqual(t, data[From[*io.Reader]()], "*io/Reader") // ?
	assertEqual(t, data[From[int]()], "int")
	assertEqual(t, data[From[int]("some-tag")], "int#some-tag")
	assertEqual(t, data[From[struct {
		r  rand.Rand  `tag:"r"`
		r2 rand2.Rand `tag:"r_2"`
	}]("some-tag")], "struct { r rand.Rand \"tag:\\\"r\\\"\"; r2 rand.Rand \"tag:\\\"r_2\\\"\" }#some-tag") //field RootID does not contain pkgPath
}

func TestFromType(t *testing.T) {
	s := struct {
		rand      rand.Rand
		rand2     rand2.Rand
		randPtr   *rand.Rand
		reader    io.Reader
		readerPtr *io.Reader
		i         int
		it        int `tag:"some-tag"`
		str       struct {
			r  rand.Rand  `tag:"r"`
			r2 rand2.Rand `tag:"r_2"`
		} `tag:"some-tag"`
	}{}

	data := testData()

	sTyp := reflect.TypeOf(s)
	for i := range sTyp.NumField() {
		field := sTyp.Field(i)
		id := FromType(field.Type, field.Tag.Get("tag"))
		_, ok := data[id]
		assertTrue(t, ok)
	}

}

func BenchmarkFrom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = From[io.Reader]("some-tag")
	}
}

func BenchmarkFromType(b *testing.B) {
	field := reflect.TypeOf(struct {
		reader io.Reader `tag:"some-tag"`
	}{}).Field(0)

	tag := field.Tag.Get("tag")

	for i := 0; i < b.N; i++ {
		_ = FromType(field.Type, tag)
	}
}
