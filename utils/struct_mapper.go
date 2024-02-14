package utils

import (
	"encoding/json"
)

// StructToMap struct to map key string
func StructToMap(in interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(in)
	json.Unmarshal(inrec, &inInterface)

	delete(inInterface, "id")
	delete(inInterface, "CreatedAt")
	delete(inInterface, "UpdatedAt")
	delete(inInterface, "DeletedAt")

	return inInterface
}
