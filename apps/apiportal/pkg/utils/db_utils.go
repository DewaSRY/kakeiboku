package utils

import (
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)


func IntToPgTypeNumeric(num int) pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(int64(num)), Valid: true}
}
