package ctxd

import (
	"encoding/json"
	"fmt"
)

// DeferredJSON postpones log field processing, suitable for debug logging.
type DeferredJSON func() interface{}

// DeferredString postpones log field processing, suitable for debug logging.
type DeferredString func() interface{}

// String implements fmt.Stringer.
func (d DeferredString) String() string {
	if d == nil {
		panic("ctxd: DeferredString is nil")
	}

	return fmt.Sprintf("%+v", d())
}

// MarshalJSON implements json.Marshaler.
func (d DeferredJSON) MarshalJSON() ([]byte, error) {
	if d == nil {
		panic("ctxd: DeferredJSON is nil")
	}

	return json.Marshal(d())
}
