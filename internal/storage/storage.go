package storage

type Storage struct {
	ShortToLong map[string]string
	LongToShort map[string]string
}

func New() *Storage {
	return &Storage{
		ShortToLong: make(map[string]string),
		LongToShort: make(map[string]string),
	}
}

func (s *Storage) AddToShort(key, value string) {
	s.ShortToLong[key] = value
}

func (s *Storage) GetFromShort(key string) (string, bool) {
	value, found := s.ShortToLong[key]
	return value, found
}

func (s *Storage) AddToLong(key, value string) {
	s.LongToShort[key] = value
}

func (s *Storage) GetFromLong(key string) (string, bool) {
	value, found := s.LongToShort[key]
	return value, found
}
