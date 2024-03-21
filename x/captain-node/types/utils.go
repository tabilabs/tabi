package types

import (
	"crypto/sha256"
	"fmt"
)

func GenDivisionsId(level int) string {
	divisionsId := fmt.Sprintf("divisions-%d", level)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(divisionsId)))
	return hash
}
