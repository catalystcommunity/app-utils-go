package test

import (
	"github.com/catalystcommunity/app-utils-go/logging"
	ulids2 "github.com/catalystcommunity/app-utils-go/ulids"
	ulid2 "github.com/oklog/ulid"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
	"time"
)

func TestUniqueOrderedUlidGeneration(t *testing.T) {
	ulids := []string{}
	ulidMap := map[string]ulid2.ULID{}
	start := time.Now()
	for time.Since(start).Milliseconds() == 0 {
		ulid := ulids2.MustNewULID()
		ulidString := ulid.String()
		_, ok := ulidMap[ulidString]
		require.False(t, ok, "there should be no duplicate ulids")
		ulidMap[ulidString] = ulid
		ulids = append(ulids, ulidString)
	}
	// create a copy of the ulids in create order
	sortedUlids := ulids
	// sort copy lexicographically
	sort.Strings(sortedUlids)
	require.Equal(t, sortedUlids, ulids, "creation order and sorted order should be the same")
	logging.Log.WithField("num_ulids", len(ulids)).Info("generated ulids")
}

func TestAccurateTimestamp(t *testing.T) {
	ulid1 := ulids2.MustNewULID()
	time.Sleep(5 * time.Millisecond)
	ulid2 := ulids2.MustNewULID()
	require.GreaterOrEqual(t, ulid2.Time()-ulid1.Time(), uint64(5), "timestamps should be at least 5ms apart")
	// might need to increase above 10ms depending on how testin gin CI goes, this is just making sure they're not wildly different
	require.LessOrEqual(t, ulid2.Time()-ulid1.Time(), uint64(10), "timestamps should be at most 10ms apart")
}
