package main

import (
	"fmt"

	"github.com/google/uuid"
)

func generateUUID() (id string) {
	id = fmt.Sprintf("%s", uuid.New())

	return
}
