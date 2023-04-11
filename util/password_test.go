package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashingPasswordOK(t *testing.T) {
	arg := Argon2Param{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	pwd := "wap12345"

	hashedPassword, err := GenerateHashFromPassword(pwd, arg)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	isMatch, err := ComparePasswordAndHashPassword(pwd, hashedPassword, arg)
	require.NoError(t, err)
	require.True(t, isMatch)
}

func TestHashingPasswordNotMatch(t *testing.T) {
	arg := Argon2Param{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	pwd := "wap12345"

	pwd2 := "wap12346"

	hashedPassword, err := GenerateHashFromPassword(pwd, arg)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	isMatch, err := ComparePasswordAndHashPassword(pwd2, hashedPassword, arg)
	require.NoError(t, err)
	require.False(t, isMatch)
}

func TestHashingPasswordFalseHashPasswordSaved(t *testing.T) {
	arg := Argon2Param{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	pwd := "wap12345"

	hashedPassword, err := GenerateHashFromPassword(pwd, arg)
	require.NoError(t, err)
	falseHashedPassword := fmt.Sprintf("%s$NewParams", hashedPassword)

	isMatch, err := ComparePasswordAndHashPassword(pwd, falseHashedPassword, arg)
	require.Error(t, err)
	require.False(t, isMatch)
}
