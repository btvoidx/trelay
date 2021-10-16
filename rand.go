package trelay

import (
	mathrand "math/rand"
	"time"
)

// seeded "math/rand"
var rand = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
