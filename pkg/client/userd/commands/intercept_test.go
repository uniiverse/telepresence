package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telepresenceio/telepresence/rpc/v2/manager"
	"github.com/telepresenceio/telepresence/v2/pkg/client/scout"
)

const multiline = `123
456`

func makeTempFile() *os.File {

	file, err := os.CreateTemp(".", "Test_interceptState_writeEnvToFileAndClose")
	if err != nil {
		panic(err)
	}
	return file
}

func Test_interceptState_writeEnvToFileAndClose(t *testing.T) {
	type fields struct {
		cmd             safeCobraCommand
		args            interceptArgs
		scout           *scout.Reporter
		connectorServer ConnectorServer
		managerClient   manager.ManagerClient
		env             map[string]string
		mountPoint      string
		localPort       uint16
		dockerPort      uint16
	}
	type args struct {
		file *os.File
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wants   string
		wantErr bool
	}{
		{name: "Simple env", args: args{file: makeTempFile()}, fields: fields{env: map[string]string{"SIMPLE": "123"}}, wants: "SIMPLE=123\n"},
		{name: "Multi-line env", args: args{file: makeTempFile()}, fields: fields{env: map[string]string{"MULTI": multiline}}, wants: "MULTI=\"123\\n456\"\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := &interceptState{
				cmd:             tt.fields.cmd,
				args:            tt.fields.args,
				scout:           tt.fields.scout,
				connectorServer: tt.fields.connectorServer,
				managerClient:   tt.fields.managerClient,
				env:             tt.fields.env,
				mountPoint:      tt.fields.mountPoint,
				localPort:       tt.fields.localPort,
				dockerPort:      tt.fields.dockerPort,
			}
			if err := is.writeEnvToFileAndClose(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("interceptState.writeEnvToFileAndClose() error = %v, wantErr %v", err, tt.wantErr)
			}
			contents, err := os.ReadFile(tt.args.file.Name())
			defer os.Remove(tt.args.file.Name())
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.wants, string(contents[:]))
		})
	}
}
