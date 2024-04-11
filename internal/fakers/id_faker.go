package fakers

import (
	"github.com/go-faker/faker/v4"
	"reflect"
)

func addIdFaker() {
	_ = faker.AddProvider("idFaker", func(v reflect.Value) (interface{}, error) {
		return faker.RandomInt(1, Count)
	})
}
