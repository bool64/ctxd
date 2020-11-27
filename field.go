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
	return fmt.Sprintf("%+v", d())
}

// MarshalJSON implements json.Marshaler.
func (d DeferredJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(d())
}
