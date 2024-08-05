package ulids

import (
	"github.com/catalystcommunity/app-utils-go/errorutils"
	"github.com/oklog/ulid"
	"math/rand"
	"sync"
	"time"
)

var Generator = NewMonotonicULIDsource(rand.New(rand.NewSource(time.Unix(1000000, 0).UnixNano())))

type MonotonicULIDsource struct {
	sync.Mutex            // mutex to allow clean concurrent access
	entropy    *rand.Rand // the entropy source
	lastMs     uint64     // the last millisecond timestamp it encountered
	lastULID   ulid.ULID  // the last ULID it generated using "github.com/oklog/ulid"
}

func NewULID() (ulid.ULID, error) {
	return Generator.New(time.Now())
}
func MustNewULID() ulid.ULID {
	ulid, err := NewULID()
	errorutils.PanicOnErr(nil, "error generating ulid", err)
	return ulid
}
func NewMonotonicULIDsource(entropy *rand.Rand) *MonotonicULIDsource {
	// get an initial ULID to kick the monotonic generation off with
	inital, err := ulid.New(ulid.Now(), entropy)
	if err != nil {
		panic(err)
	}

	return &MonotonicULIDsource{
		entropy:  entropy,
		lastMs:   0,
		lastULID: inital,
	}
}

func (u *MonotonicULIDsource) New(t time.Time) (ulid.ULID, error) {
	u.Lock()
	defer u.Unlock()

	ms := ulid.Timestamp(t)
	var err error
	if ms > u.lastMs {
		u.lastMs = ms
		u.lastULID, err = ulid.New(ms, u.entropy)
		return u.lastULID, err
	}

	// if the ms are the same, increment the entropy part of the last ULID
	// and use it as the entropy part of the new ULID.
	incrEntropy := incrementBytes(u.lastULID.Entropy())
	var dup ulid.ULID
	dup.SetTime(ms)
	if err := dup.SetEntropy(incrEntropy); err != nil {
		return dup, err
	}
	u.lastULID = dup
	u.lastMs = ms
	return dup, nil
}

func incrementBytes(in []byte) []byte {
	const (
		minByte byte = 0
		maxByte byte = 255
	)

	// copy the byte slice
	out := make([]byte, len(in))
	copy(out, in)

	// iterate over the byte slice from right to left
	// most significant byte == first byte (big-endian)
	leastSigByteIdx := len(out) - 1
	mostSigByteIdex := 0
	for i := leastSigByteIdx; i >= mostSigByteIdex; i-- {

		// If its value is 255, rollover back to 0 and try the next byte.
		if out[i] == maxByte {
			out[i] = minByte
			continue
		}
		// Else: increment.
		out[i]++
		return out
	}
	// panic if the increments are exhausted
	panic(ulid.ErrOverflow)
}
