package models

import (
	"github.com/google/uuid"
)

type Category struct {
	Id           uuid.UUID
	CategoryType bool
	Name         string
	Mandatory    bool
}
