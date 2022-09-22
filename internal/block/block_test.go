package block

import (
	"context"
	"github.com/INFURA/infra/repository/infura"
	"testing"
)

func TestGetBlock(t *testing.T) {
	repo := infura.NewInfura()
	t.Run("ByNumber", func(t *testing.T) {
		bq, err := GetBlock(context.Background(), repo, int64(6008149), "", false)
		if err != nil {
			t.Fatal(err)
		}
		if bq.Block.Number().Int64() != int64(6008149) {
			t.Errorf("expected value: %d, got: %d", 6008149, bq.Block.Number().Int64())
		}
	})
	t.Run("ByHash", func(t *testing.T) {
		bq, err := GetBlock(context.Background(), repo, 0, "0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35", false)
		if err != nil {
			t.Fatal(err)
		}
		if bq.Block.Number().Int64() != int64(6008149) {
			t.Errorf("expected value: %d, got: %d", 6008149, bq.Block.Number().Int64())
		}
	})
}
