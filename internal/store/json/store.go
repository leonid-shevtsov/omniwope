package json

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/leonid-shevtsov/omniwope/internal/store"
)

// NOTE: this provider is not safe for concurrency.
// For now, I think it's better to submit every post sequentially.
// To make it concurrent, you could have multiple files instead of one.
// (However, also for now, a single-file store is easier to transport.)
type KVProvider struct {
	path string
	data map[string]map[string]json.RawMessage
}

func (p *KVProvider) GetKV(namespace string) (store.KV, error) {
	return &KV{
		namespace: namespace,
		provider:  p,
	}, nil
}

func NewProvider(path string) store.Provider {
	return &KVProvider{path: path}
}

type KV struct {
	namespace string
	provider  *KVProvider
}

func (kv *KV) Get(key string, value interface{}) (bool, error) {
	if kv.provider.data == nil {
		err := kv.provider.readInData()
		if err != nil {
			return false, err
		}
	}

	if kv.provider.data[kv.namespace] == nil || kv.provider.data[kv.namespace][key] == nil {
		return false, nil
	}

	err := json.Unmarshal(kv.provider.data[kv.namespace][key], value)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (kv *KV) Set(key string, value interface{}) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if kv.provider.data[kv.namespace] == nil {
		kv.provider.data[kv.namespace] = make(map[string]json.RawMessage)
	}

	kv.provider.data[kv.namespace][key] = encoded

	err = kv.provider.writeData()
	if err != nil {
		return err
	}
	return nil
}

func (p *KVProvider) readInData() error {
	// read in the file
	file, err := os.Open(p.path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		p.data = make(map[string]map[string]json.RawMessage)
		return nil
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&p.data)
	if err != nil {
		return err
	}

	return nil
}

func (p *KVProvider) writeData() error {

	storeContents, err := json.Marshal(p.data)
	if err != nil {
		return err
	}
	err = os.WriteFile(p.path, storeContents, 0644)
	if err != nil {
		return err
	}
	return nil

}
