package ctxd_test

import (
	"testing"

	"github.com/bool64/ctxd"
	"github.com/stretchr/testify/assert"
)

func TestDeferredString_String(t *testing.T) {
	assert.Equal(t, "[1 2 3]", ctxd.DeferredString(func() interface{} { return []int{1, 2, 3} }).String())
}

func TestDeferredJSON_MarshalJSON(t *testing.T) {
	v, err := ctxd.DeferredJSON(func() interface{} { return []int{1, 2, 3} }).MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "[1,2,3]", string(v))
}
