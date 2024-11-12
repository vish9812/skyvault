package utils

import "github.com/google/uuid"

type ID uuid.UUID

func NewID() ID {
	id, _ := uuid.NewV7()
	return ID(id)
}

func (id ID) ToString() string {
	return uuid.UUID(id).String()
}

func (id ID) ToUUID() uuid.UUID {
	return uuid.UUID(id)
}

func ToID(uuid uuid.UUID) ID {
	return ID(uuid)
}
