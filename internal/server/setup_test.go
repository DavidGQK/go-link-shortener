package server

type TestConfig struct {
	ServerURL    string
	ShortURLBase string
}

var TestCfg = TestConfig{
	ServerURL:    "localhost:8080",
	ShortURLBase: "http://localhost:8080/",
}

type TestStorage struct {
	ShortToLong map[string]string
	LongToShort map[string]string
}

func NewTestStorage() *TestStorage {
	return &TestStorage{
		ShortToLong: make(map[string]string),
		LongToShort: make(map[string]string),
	}
}

func (s *TestStorage) AddToShort(key, value string) {
	s.ShortToLong[key] = value
}

func (s *TestStorage) GetFromShort(key string) (string, bool) {
	value, found := s.ShortToLong[key]
	return value, found
}

func (s *TestStorage) AddToLong(key, value string) {
	s.LongToShort[key] = value
}

func (s *TestStorage) GetFromLong(key string) (string, bool) {
	value, found := s.LongToShort[key]
	return value, found
}
