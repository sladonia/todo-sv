package set

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type Set struct {
	data map[string]struct{}
}

func NewSet(vals ...string) *Set {
	data := make(map[string]struct{})

	for _, val := range vals {
		data[val] = struct{}{}
	}

	return &Set{data: data}
}

func (s *Set) Contains(val string) bool {
	_, ok := s.data[val]

	return ok
}

func (s *Set) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Set) Add(vals ...string) {
	for _, val := range vals {
		s.data[val] = struct{}{}
	}
}

func (s *Set) Remove(vals ...string) {
	for _, val := range vals {
		delete(s.data, val)
	}
}

func (s *Set) IsEqual(other *Set) bool {
	return reflect.DeepEqual(s.data, other.data)
}

func (s *Set) ReplaceValues(vals ...string) {
	data := make(map[string]struct{})

	for _, val := range vals {
		data[val] = struct{}{}
	}

	s.data = data
}

func (s *Set) Values() []string {
	values := make([]string, len(s.data))

	i := 0
	for val := range s.data {
		values[i] = val
		i++
	}

	return values
}

func (s *Set) MarshalBSONValue() (bsontype.Type, []byte, error) {
	arr := bson.A{}

	i := 0
	for val := range s.data {
		arr = append(arr, val)
		i++
	}
	_, b, err := bson.MarshalValue(arr)

	return bsontype.Array, b, err
}

func (s *Set) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	var raw bson.Raw

	err := bson.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	rawNodes, err := raw.Values()
	if err != nil {
		return err
	}

	data := make(map[string]struct{}, len(rawNodes))

	for _, node := range rawNodes {
		data[node.String()] = struct{}{}
	}

	s.data = data

	return nil
}
