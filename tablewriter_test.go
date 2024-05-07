package tablewriter_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/djthorpe/go-tablewriter"
	"github.com/stretchr/testify/assert"
)

type TestAB struct {
	A string `json:"a,omitempty"`
	B string `json:"b"`
	C string `json:"-"`
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
	assert.Equal(3, meta.NumField())
	t.Log(meta)
}
