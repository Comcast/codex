package cipher

import (
	"encoding/json"
	"github.com/Comcast/webpa-common/logging"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	"github.com/spf13/viper"
)

type LocalCerts struct {
	Path     string
	HashName string
}

const (
	// CipherKey is the Viper subkey under which logging should be stored.
	// NewOptions *does not* assume this key.
	CipherKey = "cipher"
)

type Options struct {
	Logger log.Logger

	o []Config
}

type Ciphers struct {
	Logger log.Logger

	options map[string]Decrypt
}

func (o Options) LoadEncrypt() (Encrypt, error) {
	if o.Logger == nil {
		o.Logger = logging.DefaultLogger()
	}
	var lastErr error
	for _, elem := range o.o {
		if encrypter, err := elem.LoadEncrypt(); err == nil {
			elem.Logger = o.Logger
			return encrypter, nil
		} else {
			lastErr = err
		}
	}
	return DefaultCipherEncrypter(), emperror.Wrap(lastErr, "failed to load encrypt options")
}

func PopulateCiphers(o Options) Ciphers {
	c := Ciphers{
		options: map[string]Decrypt{},
	}
	if o.Logger == nil {
		o.Logger = logging.DefaultLogger()
	}
	for _, elem := range o.o {
		elem.Logger = o.Logger
		if decrypter, err := elem.LoadDecrypt(); err == nil {
			c.options[elem.KID] = decrypter
		}
	}
	return c
}

func (c *Ciphers) Get(KID string) (Decrypt, bool) {
	if d, ok := c.options[KID]; ok {
		return d, ok
	}
	return nil, false
}

// FromViper produces an Options from a (possibly nil) Viper instance.
// cipher key is expected
func FromViper(v *viper.Viper) (o Options, err error) {
	obj := v.Get("cipher")
	data, err := json.Marshal(obj)
	if err != nil {
		return Options{o: []Config{}}, emperror.Wrap(err, "failed to load cipher config")
	}

	err = json.Unmarshal(data, &o)
	return
}
