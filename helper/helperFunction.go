package helper

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetParamAsInt32(key string) (int32, error) {
	intVal, err := strconv.Atoi(key)
	return int32(intVal), err
}

// Convert pgtype.Numeric to string
func NumericToString(n pgtype.Numeric) (string, error) {
	if !n.Valid || n.Int == nil {
		return "", fmt.Errorf("invalid numeric value")
	}

	rat := new(big.Rat).SetInt(n.Int)

	// Apply the exponent: 10^(-n.Exp)
	scale := big.NewInt(1)
	scale.Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil)

	rat.Quo(rat, new(big.Rat).SetInt(scale))

	return rat.FloatString(max(0, int(-n.Exp))), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
