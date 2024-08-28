package callback

import (
	"strconv"
	"strings"
)

func GetID(data string) int {
	parts := strings.Split(data, "_")
	if len(parts) > 3 {
		return 0
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0
	}

	return id
}
