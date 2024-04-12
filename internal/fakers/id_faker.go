package fakers

import (
	"github.com/go-faker/faker/v4"
	"math/rand/v2"
	"reflect"
)

func addIdFaker() {
	_ = faker.AddProvider("idFaker", func(v reflect.Value) (interface{}, error) {
		return rand.IntN(Count) + 1, nil
	})
}
