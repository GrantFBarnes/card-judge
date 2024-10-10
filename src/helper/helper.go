package helper

import "github.com/google/uuid"

func ConvertArrayStringsToUuids(stringArray []string) []uuid.UUID {
	uuidArray := make([]uuid.UUID, len(stringArray))
	for i := range stringArray {
		id, err := uuid.Parse(stringArray[i])
		if err != nil {
			uuidArray[i] = uuid.Nil
			continue
		}
		uuidArray[i] = id
	}
	return uuidArray
}

func ConvertArrayUuidsToStrings(uuidArray []uuid.UUID) []string {
	stringArray := make([]string, len(uuidArray))
	for i := range uuidArray {
		stringArray[i] = uuidArray[i].String()
	}
	return stringArray
}

func IsIdInArray(idToFind uuid.UUID, uuidArray []uuid.UUID) bool {
	for _, id := range uuidArray {
		if id == idToFind {
			return true
		}
	}
	return false
}
