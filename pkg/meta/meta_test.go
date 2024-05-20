package meta_test

import (
	"reflect"
	"testing"

	meta "github.com/djthorpe/go-tablewriter/pkg/meta"
	assert "github.com/stretchr/testify/assert"
)

type TestAB struct {
	A string `json:"a,omitempty" writer:",test1:test2"`
	B string `json:"b" writer:"-"`
	C string `json:"-"`
}

type TestABEF struct {
	TestAB
	E string
	F string
}

type TestG struct {
	G *string `description:"this is field G"`
}

func Test_meta_000(t *testing.T) {
	assert := assert.New(t)
	meta, err := meta.New(TestAB{})
	if !assert.NoError(err) {
		t.SkipNow()
	}
	assert.NotNil(meta)
}

func Test_meta_002(t *testing.T) {
	assert := assert.New(t)
	meta, err := meta.New(TestAB{})
	if !assert.NoError(err) {
		t.SkipNow()
	}
	assert.NotNil(meta)
	assert.Equal(reflect.TypeOf(TestAB{}), meta.Type())
}

func Test_meta_003(t *testing.T) {
	assert := assert.New(t)
	meta, err := meta.New([]TestAB{}, "json")
	assert.NoError(err)
	assert.NotNil(meta)
	assert.Equal(reflect.TypeOf(TestAB{}), meta.Type())

	fields := meta.Fields()
	assert.Equal(2, len(fields)) // A and B
	t.Log(meta)
}

func Test_meta_004(t *testing.T) {
	assert := assert.New(t)
	meta, err := meta.New([]TestABEF{}, "json")
	assert.NoError(err)
	assert.NotNil(meta)
	assert.Equal(reflect.TypeOf(TestABEF{}), meta.Type())

	fields := meta.Fields()
	assert.Equal(4, len(fields)) // A , B, E , F
	t.Log(meta)
}

func Test_meta_005(t *testing.T) {
	assert := assert.New(t)
	meta, err := meta.New([]TestAB{}, "writer", "json")
	assert.NoError(err)
	assert.NotNil(meta)
	assert.Equal(reflect.TypeOf(TestAB{}), meta.Type())

	fields := meta.Fields()
	assert.Equal(1, len(fields)) // A
	assert.Equal("a", fields[0].Name())
	assert.True(fields[0].Is("omitempty"))
	assert.Equal("test2", fields[0].Tuple("test1"))
	t.Log(meta)
}

func Test_meta_006(t *testing.T) {
	assert := assert.New(t)
	meta, err := meta.New(&TestG{})
	assert.NoError(err)
	assert.NotNil(meta)
	assert.Equal(reflect.TypeOf(TestG{}), meta.Type())

	fields := meta.Fields()
	assert.Equal(1, len(fields)) // G
	assert.Equal("this is field G", fields[0].Tag("description"))
}
