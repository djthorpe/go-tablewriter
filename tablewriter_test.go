package tablewriter_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/djthorpe/go-tablewriter"
	"github.com/stretchr/testify/assert"
)

type TestAB struct {
	A string `json:"a,omitempty"`
	B string `json:"b"`
	C string `json:"-"`
}

type TestABEF struct {
	TestAB
	E string
	F string
}

func Test_tablewriter_000(t *testing.T) {
	assert := assert.New(t)
	writer := tablewriter.New(os.Stdout)
	assert.NotNil(writer)
}

func Test_tablewriter_001(t *testing.T) {
	assert := assert.New(t)
	writer := tablewriter.New(os.Stdout)
	err := writer.Write([]struct{}{})
	assert.NoError(err)
}

func Test_tablewriter_002(t *testing.T) {
	assert := assert.New(t)
	writer := tablewriter.New(os.Stdout)
	meta, err := writer.NewMeta(TestAB{})
	assert.NoError(err)
	assert.NotNil(meta)
	assert.Equal(reflect.TypeOf(TestAB{}), meta.Type())
}

func Test_tablewriter_003(t *testing.T) {
	assert := assert.New(t)
	writer := tablewriter.New(os.Stdout)
	meta, err := writer.NewMeta([]TestAB{})
	assert.NoError(err)
	assert.NotNil(meta)
	assert.Equal(reflect.TypeOf(TestAB{}), meta.Type())
	assert.Equal(2, meta.NumField()) // A and B
	assert.Equal([]string{"a", "b"}, meta.Fields())
	t.Log(meta)
}

func Test_tablewriter_004(t *testing.T) {
	assert := assert.New(t)
	writer := tablewriter.New(os.Stdout)
	meta, err := writer.NewMeta([]TestABEF{})
	assert.NoError(err)
	assert.NotNil(meta)
	assert.Equal(reflect.TypeOf(TestABEF{}), meta.Type())
	assert.Equal(4, meta.NumField()) // A,B,E,F
	assert.Equal([]string{"a", "b", "E", "F"}, meta.Fields())
	t.Log(meta)
}

func Test_tablewriter_005(t *testing.T) {
	assert := assert.New(t)
	buf := new(strings.Builder)
	writer := tablewriter.New(buf)
	table := []TestAB{
		{A: "hello", B: "world"},
	}
	err := writer.Write(table)
	assert.NoError(err)
	assert.Equal("hello,world\n", buf.String())
}

func Test_tablewriter_006(t *testing.T) {
	assert := assert.New(t)
	buf := new(strings.Builder)
	writer := tablewriter.New(buf, tablewriter.OptFieldDelim('|'))
	table := []TestAB{
		{A: "hello", B: "world"},
	}
	err := writer.Write(table)
	assert.NoError(err)
	assert.Equal("hello|world\n", buf.String())
}

func Test_tablewriter_007(t *testing.T) {
	assert := assert.New(t)
	buf := new(strings.Builder)
	writer := tablewriter.New(buf, tablewriter.OptHeader())
	table := []TestAB{
		{A: "hello", B: "world"},
	}
	err := writer.Write(table)
	assert.NoError(err)
	assert.Equal("a,b\nhello,world\n", buf.String())
}
