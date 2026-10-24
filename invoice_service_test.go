package main

import (
	"fmt"
	"testing"
)

func TestInvoiceAmountDecision(t *testing.T) {
	cases := []struct {
		name, amount string
		want         bool
	}{{"positive", "125.00", true}, {"zero", "0", false}, {"invalid", "oops", false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var n float64
			if _, err := fmt.Sscan(tc.amount, &n); err != nil {
				n = 0
			}
			got := n > 0
			if got != tc.want {
				t.Fatalf("amount %q: got %v want %v", tc.amount, got, tc.want)
			}
		})
	}
}
