package fts

import (
	"testing"
)

func TestStrategyConfig(t *testing.T) {
	expected :=
		struct {
			BuyThreshold        float32
			SellThreshold       float32
			StopLossThreshold   float32
			MissProfitThreshold float32
		}{
			BuyThreshold:        12.34,
			SellThreshold:       23.45,
			StopLossThreshold:   34.56,
			MissProfitThreshold: 45.67,
		}
	gotten := get_fts_config(expected)

	if gotten.BuyThreshold != expected.BuyThreshold {
		t.Fatalf("buy thr: expected = %v, gotten = %v",
			expected.BuyThreshold, gotten.BuyThreshold)
	}
	if gotten.SellThreshold != expected.SellThreshold {
		t.Fatalf("sell thr: expected = %v, gotten = %v",
			expected.SellThreshold, gotten.SellThreshold)
	}
	if gotten.StopLossThreshold != expected.StopLossThreshold {
		t.Fatalf("stop loss thr: expected = %v, gotten = %v",
			expected.StopLossThreshold, gotten.StopLossThreshold)
	}
	if gotten.MissProfitThreshold != expected.MissProfitThreshold {
		t.Fatalf("miss profit thr: expected = %v, gotten = %v",
			expected.MissProfitThreshold, gotten.MissProfitThreshold)
	}
}
