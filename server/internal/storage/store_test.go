package main

import (
	"testing"
)
//	"github.com/go-errors/errors"
// 	"github.com/stretchr/testify/require"

func TestPut(t *testing.T) {
/*
	const key = "put_key"
	const value = "put_value"

	var val interface{}
	var ok bool

	defer delete(store.kvLog, key)

	// check for existing key
	_, ok = store.kvLog[key]
	require.False(t, ok, "key/value pair already exist")

	// store the value at key
	err := Put(key, value)
	require.NoError(t, err)

	val, ok = store.kvLog[key]
	require.True(t, ok, "Put new key/value failed")
	require.Equal(t, val, value)
*/
}

func TestGet(t *testing.T) {
/*
	const key = "get_key"
	const value = "get_value"

	var val interface{}
	var err error

	defer delete(store.kvLog, key)

	// check key isn't already present
	val, err = Get(key)
	require.Error(t, err)
	require.EqualError(t, err, "key not found")

	store.kvLog[key] = value

	val, err = Get(key)
	require.NoError(t, err)
	require.Equal(t, val, value)
*/
}

func TestDelete(t *testing.T) {
/*
	const key = "delete_key"
	const value = "delete_value"

	var ok bool

	defer delete(store.kvLog, key)

	store.kvLog[key] = value

	_, ok = store.kvLog[key]
	require.True(t, ok, "key/value pair doesn't exist")

	Delete(key)

	_, ok = store.kvLog[key]
	require.False(t, ok, "Delete failed")
*/
}