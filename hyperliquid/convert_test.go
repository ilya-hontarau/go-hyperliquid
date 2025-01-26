package hyperliquid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	size := RoundOrderSize(1.0001, 3)
	assert.Equal(t, "1.000", size)

	size = RoundOrderSize(1.2345, 3)
	assert.Equal(t, "1.235", size)
	size = RoundOrderSize(1.2344, 3)
	assert.Equal(t, "1.234", size)
	size = RoundOrderSize(1.2344, 0)
	assert.Equal(t, "1", size)
	size = RoundOrderSize(1.554, 0)
	assert.Equal(t, "2", size)
}

func TestPriceConvert(t *testing.T) {
	assert.Equal(t, "123457", RoundOrderPrice(123456.534, 3, 8))
	assert.Equal(t, "12456", RoundOrderPrice(12456.3, 3, 8))
	assert.Equal(t, "12456", RoundOrderPrice(12456.0, 3, 8))
	assert.Equal(t, "1234.5", RoundOrderPrice(1234.5, 3, 8))
	assert.Equal(t, "1234.5", RoundOrderPrice(1234.54, 3, 8))
	assert.Equal(t, "1234.6", RoundOrderPrice(1234.55, 3, 8))
	assert.Equal(t, "34.56", RoundOrderPrice(34.556, 6, 8))
	assert.Equal(t, "0.0001234", RoundOrderPrice(0.0001234, 0, 8))
	assert.Equal(t, "0.0001235", RoundOrderPrice(0.0001235, 0, 8))
	assert.Equal(t, "0.000123", RoundOrderPrice(0.0001234, 2, 8))
	assert.Equal(t, "0.000124", RoundOrderPrice(0.0001235, 2, 8))

	t.Run("Perps - Valid integer", func(t *testing.T) {
		assert.Equal(t, "123456", RoundOrderPrice(123456.0, 3, 8))
	})
	t.Run("Perps - Valid decimal", func(t *testing.T) {
		assert.Equal(t, "1234.5", RoundOrderPrice(1234.5, 3, 6))
	})
	t.Run("Perps - Valid decimal with trailing zeros", func(t *testing.T) {
		assert.Equal(t, "0.001", RoundOrderPrice(0.001234, 3, 6))
	})
}
