package ipgen

import (
	"fmt"
	"math/rand"
)

// Return a random IP from 203.0.113.0/24 (RFC 5737 test-net-3).
func GenTestNet3() string {
	host := rand.Intn(254) + 1
	return fmt.Sprintf("203.0.113.%d", host)
}
