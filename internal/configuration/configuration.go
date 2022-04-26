package configuration

import (
	"errors"
	"os"

	"github.com/pja237/prom2tower/internal/pipe"
	"gopkg.in/yaml.v2"
)

type config struct {
	Globals map[string]interface{}
	Glue    []pipe.Pipe
}

func WantString(v interface{}) (*string, error) {
	if s, ok := v.(string); ok {
		return &s, nil
	} else {
		return nil, errors.New("not a string")
	}
}

func ReturnString(v interface{}) *string {
	if s, ok := v.(string); ok {
		return &s
	} else {
		return nil
	}
}

func GetConfig(path string) (*config, error) {

	// open and read file into []byte
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// call unmarshal
	conf, err := unmarshalConfig(data)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func unmarshalConfig(data []byte) (*config, error) {
	var conf = new(config)

	//err := json.Unmarshal(data, conf)
	err := yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
