package pregen

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate_Default(t *testing.T) {
	t.Parallel()

	type Data struct {
		otp, hash string
	}

	generator, terminate := NewGenerator(func() (Data, error) {
		otp := "293940" // your random otp

		// hash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
		time.Sleep(80 * time.Millisecond) // mock bcrypt hashing work
		hash := []byte("foo")

		return Data{
			otp:  otp,
			hash: string(hash),
		}, nil
	},
		PregenSize[Data](19),
		ErrorCooldown[Data](100*time.Millisecond),
		StartDelay[Data](10*time.Millisecond),
	)
	defer terminate()

	time.Sleep(time.Second)

	for range 10 {
		start := time.Now()
		got, err := generator.Gen()
		assert.NotEmpty(t, got)
		require.NoError(t, err)

		assert.WithinDuration(t, start, time.Now(), time.Millisecond)
	}
}
