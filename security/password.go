package security

import (
	"errors"
	"fold/console"
	"fold/mem"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

var (
	PasswordEncoder mem.MapColumnBatch
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func EncodePasswords(table *mem.Table) {
	inputCol := table.ColNumber("input_password")
	outputCol := table.ColNumber("password")

	if inputCol < 0 || outputCol < 0 {
		console.RedPrintln("password columns not found")
		return
	}

	PasswordEncoder = mem.MapColumnBatch{
		InputColumn:  inputCol,
		ResultColumn: outputCol,
		Transform: func(s string) (string, error) {
			if strings.TrimSpace(s) == "" {
				return "", errors.New("empty password")
			}
			res, err := HashPassword(s)
			return res, err
		},
	}

	table.BatchUpdate(PasswordEncoder)
}
