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

func StringToPgText(value string) pgtype.Text {
    return pgtype.Text{
    	String: value,
    	Valid:  true,
    }
}
