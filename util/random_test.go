package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenRandomCurrent(t *testing.T) {
	for i := 0; i < 10; i++ {
		cur := GenRandomCurrency()
		require.Contains(t, "EUR USD IDR", cur)
	}
}

func TestGenRandomOwner(t *testing.T) {
	testCase := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "Check0Length",
			check: func(t *testing.T) {
				s := GenRandomOwner()
				l := len(s)
				require.Greater(t, l, 0)
			},
		},
		{
			name: "Check8length",
			check: func(t *testing.T) {
				s := GenRandomOwner()
				l := len(s)
				require.LessOrEqual(t, l, 8)
			},
		},
	}

	for _, ts := range testCase {
		t.Run(ts.name, ts.check)
	}
}

func TestGenRandomMoney(t *testing.T) {
	testCase := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "NotLessThan0",
			check: func(t *testing.T) {
				s := GenRandomMoney()
				require.Greater(t, s, int64(0))
			},
		},
		{
			name: "NotMoreThan1000000",
			check: func(t *testing.T) {
				s := GenRandomMoney()
				require.LessOrEqual(t, s, int64(1000000))
			},
		},
	}

	for _, ts := range testCase {
		t.Run(ts.name, ts.check)
	}
}
