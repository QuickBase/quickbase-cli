package qbclient_test

import (
	"testing"

	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/spf13/viper"
)

func TestInvalidInput(t *testing.T) {
	input := &qbclient.GetVariableInput{}

	cfg := qbclient.NewConfig(viper.New())
	client := qbclient.New(cfg)

	output := &qbclient.GetVariableOutput{}
	err := client.Do(input, output)

	if err == nil {
		t.Fatal("got nil, expected error")
	}
}
