package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseId(t *testing.T) {
	t.Run("valid parse", func(t *testing.T) {
		//arrange
		inputDataUint := uint64(12)
		inputData := "12"

		//act
		res, err := ParseId(inputData)

		//assert
		require.NoError(t, err)
		require.Equal(t, res, inputDataUint)
	})

	t.Run("invalid parse in long data", func(t *testing.T) {
		//arrange
		inputData := "12 11 10"

		//act
		_, err := ParseId(inputData)

		//assert
		require.Error(t, err)
		require.EqualError(t, ErrInvalidIdParser, err.Error())
	})

	t.Run("invalid parse", func(t *testing.T) {
		//arrange
		inputData := ""

		//act
		_, err := ParseId(inputData)

		//assert
		require.Error(t, err)
		require.EqualError(t, ErrInvalidIdParser, err.Error())
	})

	t.Run("invalid value", func(t *testing.T) {
		//arrange
		inputData := "aq"

		//act
		_, err := ParseId(inputData)

		//assert
		require.Error(t, err)
		require.EqualError(t, ErrInvalidIdParser, err.Error())
	})
}
