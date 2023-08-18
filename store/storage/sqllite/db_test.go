package sqllite

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	storeKey1 = "store1"
)

func TestFoo(t *testing.T) {
	db, err := New(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, db.Close())
}

func TestDatabase_LatestVersion(t *testing.T) {
	db, err := New(t.TempDir())
	require.NoError(t, err)
	defer db.Close()

	lv, err := db.GetLatestVersion()
	require.Error(t, err)
	require.Zero(t, lv)

	for i := uint64(1); i <= 100; i++ {
		err = db.SetLatestVersion(i)
		require.NoErrorf(t, err, "failed to set latest version: %d", i)

		lv, err = db.GetLatestVersion()
		require.NoErrorf(t, err, "failed to get latest version: %d", i)
		require.Equal(t, i, lv)
	}
}

func TestDatabase_CRUD(t *testing.T) {
	db, err := New(t.TempDir())
	require.NoError(t, err)
	defer db.Close()

	ok, err := db.Has(storeKey1, 1, []byte("key"))
	require.NoError(t, err)
	require.False(t, ok)

	err = db.Set(storeKey1, 1, []byte("key"), []byte("value"))
	require.NoError(t, err)

	ok, err = db.Has(storeKey1, 1, []byte("key"))
	require.NoError(t, err)
	require.True(t, ok)

	val, err := db.Get(storeKey1, 1, []byte("key"))
	require.NoError(t, err)
	require.Equal(t, []byte("value"), val)

	err = db.Delete(storeKey1, 1, []byte("key"))
	require.NoError(t, err)

	ok, err = db.Has(storeKey1, 1, []byte("key"))
	require.NoError(t, err)
	require.False(t, ok)

	val, err = db.Get(storeKey1, 1, []byte("key"))
	require.NoError(t, err)
	require.Nil(t, val)

	err = db.Delete(storeKey1, 1, []byte("not_exists"))
	require.NoError(t, err)
}
