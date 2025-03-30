package crud

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("value for key not found")
	ErrEmptyKey = errors.New("validation of key failed, key can't be empty")
)

type (
	Service interface {
		Create(cx context.Context, key string, value string) error
		Get(cx context.Context, key string) (*KeyValue, error)
		Delete(cx context.Context, key string) error
		List(cx context.Context) (keys []string, err error)
	}

	//

	KeyValue struct {
		Key   string
		Value string
	}

	defaultCRUDService struct {
		store map[string]*KeyValue
	}
)

func NewKeyValue(key, value string) *KeyValue {
	return &KeyValue{key, value}
}

func (dns *defaultCRUDService) Create(
	cx context.Context, key, value string,
) (err error) {
	if key == "" {
		return ErrEmptyKey
	}

	dns.store[key] = NewKeyValue(key, value)
	return err
}

func (dns *defaultCRUDService) Get(cx context.Context, key string) (*KeyValue, error) {
	if _, ok := dns.store[key]; !ok {
		return nil, ErrNotFound
	}
	return dns.store[key], nil
}

func (dns *defaultCRUDService) Delete(cx context.Context, key string) (err error) {
	delete(dns.store, key)
	return err
}

func (dns *defaultCRUDService) List(cx context.Context) (keys []string, err error) {
	keys = make([]string, 0)
	for k := range dns.store {
		keys = append(keys, k)
	}
	return keys, err

}

func NewCrudService(samplekey, samplevalue string) (Service, error) {
	dcs := &defaultCRUDService{
		store: make(map[string]*KeyValue),
	}

	dcs.store[samplekey] = NewKeyValue(samplekey, samplevalue)
	return dcs, nil
}
