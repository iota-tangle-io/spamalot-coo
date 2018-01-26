package server

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"
)

type Config interface{}

var subconfigs = []Config{&AppConfig{}, &NetConfig{}}

func LoadConfig() *Configuration {
	configuration := &Configuration{}
	refConfig := reflect.Indirect(reflect.ValueOf(configuration))

	// go through each sub config, load it and init it on the main struct
	for _, c := range subconfigs {
		// indirect as 'c' is pointer to struct
		ind := reflect.Indirect(reflect.ValueOf(c))
		ty := ind.Type()
		field, _ := ty.FieldByName("Location")
		fileLocation := field.Tag.Get("loc")

		// read file indicated by the field tag
		fileBytes, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(fileBytes, c); err != nil {
			panic(err)
		}

		// init configuration struct field with the given config
		configFieldName := strings.Split(ty.Name(), "Config")[0]
		refConfig.FieldByName(configFieldName).Set(ind)
	}
	return configuration
}

type Configuration struct {
	App AppConfig
	Net NetConfig
}

type AppConfig struct {
	Location interface{} `loc:"./configs/app.json"`
	Name     string
	Dev      bool
	Verbose  bool
}

type NetConfig struct {
	Location     interface{} `loc:"./configs/network.json"`
	HTTP         WebConfig
	LDAP         LDAPConfig
	Database     DatabaseConfig
	SessionStore SessionStoreConfig
}

type WebConfig struct {
	Domain  string
	Address string
	TLS     struct {
		Use         bool
		PrivateKey  string
		Certificate string
	}
	Assets struct {
		Static  string
		HTML    string
		Favicon string
	}
	ReCaptcha struct {
		Use        bool
		PrivateKey string
		PublicKey  string
	}
	Sessions struct {
		PrivateKey string
	}
	LogRequests bool
}

type LDAPConfig struct {
	Address  string
	Port     int
	User     string
	Password string
	Domain   string
	Base     string
	GroupOU  string
	UserOUs  []string
}

type DatabaseConfig struct {
	Mongo MongoDBConfig
}

type MongoDBConfig struct {
	Use       bool
	Address   string
	Auth      bool
	Username  string
	Password  string
	Mechanism string
	Source    string
}

type SessionStoreConfig struct {
	Redis RedisConfig
}

type RedisConfig struct {
	Address string
}
