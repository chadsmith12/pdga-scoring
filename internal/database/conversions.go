package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func IntToPgInt(value int) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(value),
		Valid: true,
	}
}

func BigIntToPgInt8(value int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64: int64(value),
		Valid: true,
	}
}

func StringToPgText(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}

func BoolToPgBool(value bool) pgtype.Bool {
	return pgtype.Bool{
		Bool:  value,
		Valid: true,
	}
}
