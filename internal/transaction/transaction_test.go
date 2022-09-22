package transaction

import (
	"context"
	"github.com/INFURA/infra/repository/infura"
	"testing"
)

func TestGetTransaction(t *testing.T) {
	t.Run("ByNumberAndIndex", func(t *testing.T) {
		tx, err := GetTransaction(context.Background(), infura.NewInfura(), "", int64(6008149), 0, "")
		if err != nil {
			t.Fatal(err)
		}
		if tx.Hash != "0x8784d99762bccd03b2086eabccee0d77f14d05463281e121a62abfebcf0d2d5f" {
			t.Errorf("expected valued: %s, got: %s", "0x8784d99762bccd03b2086eabccee0d77f14d05463281e121a62abfebcf0d2d5f", tx.Hash)
		}
	})

	t.Run("ByNumberAndIndex-Latest", func(t *testing.T) {
		tx, err := GetTransaction(context.Background(), infura.NewInfura(), "latest", 0, 0, "")
		if err != nil {
			t.Fatal(err)
		}
		if tx.Gas == "" {
			t.Errorf("unexpected value")
		}
	})
	t.Run("ByHash", func(t *testing.T) {
		tx, err := GetTransaction(context.Background(), infura.NewInfura(), "", 0, -1, "0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0")
		if err != nil {
			t.Fatal(err)
		}
		if tx.Hash != "0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0" {
			t.Errorf("expected value: %s, got: %s", "0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0", tx.Hash)
		}
	})

	t.Run("ByHashAndIndex", func(t *testing.T) {
		tx, err := GetTransaction(context.Background(), infura.NewInfura(), "", 0, 0, "0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35")
		if err != nil {
			t.Fatal(err)
		}
		if tx.Hash != "0x8784d99762bccd03b2086eabccee0d77f14d05463281e121a62abfebcf0d2d5f" {
			t.Errorf("expected value: %s, got: %s", "0x8784d99762bccd03b2086eabccee0d77f14d05463281e121a62abfebcf0d2d5f", tx.Hash)
		}
	})
}
