package tablewriter_test

import (
	"os"
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

type TestG struct {
	G *string
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

func Test_tablewriter_005(t *testing.T) {
	assert := assert.New(t)
	buf := new(strings.Builder)
	writer := tablewriter.New(buf)
	table := []TestAB{
		{A: "hello", B: "world"},
	}
	err := writer.Write(table, tablewriter.OptOutputCSV())
	assert.NoError(err)
	assert.Equal("hello,world\n", buf.String())
}

func Test_tablewriter_006(t *testing.T) {
	assert := assert.New(t)
	buf := new(strings.Builder)
	writer := tablewriter.New(buf, tablewriter.OptOutputCSV())
	table := []TestAB{
		{A: "hello", B: "world"},
	}
	err := writer.Write(table, tablewriter.OptDelimiter('|'))
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

func Test_tablewriter_008(t *testing.T) {
	assert := assert.New(t)
	buf := new(strings.Builder)
	writer := tablewriter.New(buf, tablewriter.OptNull("NULL"))
	table := TestG{}
	err := writer.Write(table)
	assert.NoError(err)
	assert.Equal("NULL\n", buf.String())
}
