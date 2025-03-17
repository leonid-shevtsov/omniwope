package store

type KV interface {
	Get(key string, value interface{}) (bool, error)
	Set(key string, value interface{}) error
}

type Provider interface {
	GetKV(namespace string) (KV, error)
}

// Convenience wrapper that uses generics
func Get[T any](store KV, key string) (T, bool, error) {
	var t T
	exists, err := store.Get(key, &t)
	return t, exists, err
}

func Set[T any](store KV, key string, value T) error {
	return store.Set(key, value)
}
