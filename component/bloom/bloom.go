package bloom

import (
	"github.com/go-redis/redis/v7"
)

type Filter struct {
	HashFunc []func([]byte) uint64
	bitSet   bitSetProvider
	bits     uint
}

type bitSetProvider interface {
	check([]uint) (bool, error)
	set([]uint) error
}

type bitSet struct {
	store *redis.Client
	key   string
	bits  uint
}

const (
	setScript = `
for _, offset in ipairs(ARGV) do
	redis.call("setbit", KEYS[1], offset, 1)
end
`
	checkScript = `
for _,offset in ipairs(ARGV) do
	if tonumber(redis.call("gitbit", KEYS[1],offset)) == 0 then
		return false
   end
end
return true
`
)

func New(store *redis.Client, key string, bits uint, funcs []func([]byte) uint64) *Filter {
	return &Filter{
		HashFunc: funcs,
		bitSet: &bitSet{
			store: store,
			key:   key,
			bits:  bits,
		},
	}
}

// Add adds data into f.
func (f *Filter) Add(data []byte) error {
	locations := f.getLocations(data)
	return f.bitSet.set(locations)
}

//
// Exists checks if data is in f.
func (f *Filter) Exists(data []byte) (bool, error) {
	locations := f.getLocations(data)
	isSet, err := f.bitSet.check(locations)
	if err != nil {
		return false, err
	}
	if !isSet {
		return false, nil
	}

	return true, nil
}

func (f *Filter) getLocations(data []byte) []uint {
	locations := make([]uint, len(f.HashFunc))
	for i := 0; i < len(f.HashFunc); i++ {
		hashValue := f.HashFunc[i](data)
		locations[i] = uint(hashValue % uint64(f.bits))
	}
	return locations
}

func (s *bitSet) set(offsets []uint) error {
	err := s.store.Eval(setScript, []string{s.key}, offsets).Err()
	if err == redis.Nil {
		return nil
	}
	return err
}

func (s *bitSet) check(offsets []uint) (bool, error) {
	val := s.store.Eval(checkScript, []string{s.key}, offsets).Val()

	exists, ok := val.(int64)
	if !ok {
		return false, nil
	}
	return exists == 1, nil
}
