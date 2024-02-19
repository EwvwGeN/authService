package verification

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HappyPass(t *testing.T) {
	uuid := gofakeit.UUID()

	vCode, err := GenerateVerificationCode([]byte(uuid))
	require.NoError(t, err)

	decodedUid, err := DecryptVerificationCode(vCode)
	require.NoError(t, err)

	assert.Equal(t, string(decodedUid), uuid)
}

func Test_Fall(t *testing.T) {
	uuid := gofakeit.UUID()

	vCode, err := GenerateVerificationCode([]byte(uuid))
	require.NoError(t, err)

	decodedUid, err := DecryptVerificationCode(vCode)
	require.NoError(t, err)

	uuid = gofakeit.UUID()
	assert.NotEqual(t, string(decodedUid), uuid)
}

func Test_OneKeyForSession(t *testing.T) {
	uuid := gofakeit.UUID()

	vCode, err := GenerateVerificationCode([]byte(uuid))
	require.NoError(t, err)

	decodedUid, err := DecryptVerificationCode(vCode)
	require.NoError(t, err)

	require.Equal(t, string(decodedUid), uuid)

	_, err = GenerateVerificationCode([]byte(uuid))
	require.NoError(t, err)

	decodedUid, err = DecryptVerificationCode(vCode)
	require.NoError(t, err)

	require.Equal(t, string(decodedUid), uuid)
}