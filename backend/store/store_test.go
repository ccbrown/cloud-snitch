package store_test

import (
	"testing"

	"github.com/ccbrown/cloud-snitch/backend/store"
	"github.com/ccbrown/cloud-snitch/backend/store/storetest"
)

func NewTestStore(t *testing.T) *store.Store {
	config := storetest.NewStoreConfig(t)
	ret, err := store.New(config)
	if err != nil {
		t.Fatal(err)
	}
	return ret
}
