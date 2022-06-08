package set

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMarshalling(t *testing.T) {
	type Wrapper struct {
		Array *Set `bson:"array"`
	}

	s := NewSet()
	s.Add("first", "second")

	dummy := Wrapper{Array: s}

	b, err := bson.Marshal(dummy)
	assert.NoError(t, err)

	marshalExample()
	fmt.Println(string(b))

	var result Wrapper

	err = bson.Unmarshal(b, &result)
	assert.NoError(t, err)

	fmt.Printf("%+v\n", result.Array)
}

func marshalExample() {
	type Wrapper struct {
		Array []string `bson:"array"`
	}

	dummy := Wrapper{Array: []string{"first", "second"}}

	b, err := bson.Marshal(dummy)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
