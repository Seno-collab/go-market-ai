package db

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/jackc/pgtype"
)

func NumericToFloat(n pgtype.Numeric) (float64, error) {
	if n.NaN {
		return 0, fmt.Errorf("numeric is NaN")
	}

	if n.Int == nil {
		return 0, fmt.Errorf("numeric has no value (nil Int)")
	}

	// Convert big.Int to float64
	f, _ := new(big.Float).SetInt(n.Int).Float64()

	// Apply exponent
	f = f * math.Pow10(int(n.Exp))

	return f, nil
}

func NumericToInt(n pgtype.Numeric) (int64, error) {
	f, err := NumericToFloat(n)
	return int64(f), err
}

func NumericToString(n pgtype.Numeric) (string, error) {
	if n.NaN {
		return "", fmt.Errorf("numeric is NaN")
	}
	if n.Int == nil {
		return "", fmt.Errorf("nil numeric")
	}

	s := n.Int.String()

	if n.Exp < 0 {
		scale := int(-n.Exp)

		if len(s) <= scale {
			zeros := strings.Repeat("0", scale-len(s)+1)
			s = zeros + s
		}

		idx := len(s) - scale
		s = s[:idx] + "." + s[idx:]
	}

	return s, nil
}
