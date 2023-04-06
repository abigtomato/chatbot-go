package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := LoadConfigWithPath("../config.dev.yaml")
	logrus.Infoln(fmt.Sprintf("%+v", config))
}
