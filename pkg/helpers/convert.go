package helpers

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
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

func NumericToMoney(n pgtype.Numeric) (Money, error) {
	if !n.Valid {
		return 0, nil
	}

	if n.Int == nil {
		return 0, fmt.Errorf("numeric has no value (nil Int)")
	}

	v := new(big.Int).Set(n.Int)
	switch {
	case n.Exp > 0:
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n.Exp)), nil)
		v.Mul(v, scale)
	case n.Exp < 0:
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil)
		remainder := new(big.Int).Mod(v, scale)
		if remainder.Sign() != 0 {
			return 0, fmt.Errorf("numeric has decimal part, cannot convert to Money")
		}
		v.Quo(v, scale)
	}

	if !v.IsInt64() {
		return 0, fmt.Errorf("numeric overflows int64")
	}

	return NewMoney(v.Int64())
}

func NumericFromMoney(m Money) pgtype.Numeric {
	return m.Numeric()
}

func (m Money) Numeric() pgtype.Numeric {
	return pgtype.Numeric{
		Int:   big.NewInt(int64(m)),
		Exp:   0, // integer, 10^0
		Valid: true,
	}
}
