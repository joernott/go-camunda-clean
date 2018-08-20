package camunda

import (
	"encoding/json"

	"github.com/joernott/lra"

	_ "github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

type CamundaConnection struct {
	Connection *lra.Connection
}

func NewCamunda(UseSSL bool, Server string, Port int, BaseEndpoint string, User string, Password string, ValidateSSL bool, Proxy string, ProxyIsSocks bool) (*CamundaConnection, error) {
	cam := new(CamundaConnection)
	SendHeaders := make(lra.HeaderList)
	SendHeaders["Content-Type"] = "application/json"
	c, err := lra.NewConnection(UseSSL, Server, Port, BaseEndpoint, User, Password, ValidateSSL, Proxy, ProxyIsSocks, SendHeaders)
	if err != nil {
		return nil, err
	}
	cam.Connection = c
	return cam, nil
}

func (camunda *CamundaConnection) GetProcessInstanceList() (ProcessInstanceList, error) {
	var list ProcessInstanceList

	logger := log.WithField("func", "GetProcessInstanceList").Logger
	raw, err := camunda.Connection.Get("/process-instance")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	err = json.Unmarshal(raw, &list)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return list, nil
}

func (camunda *CamundaConnection) TerminateProcess(Id string) error {

	logger := log.WithField("func", "GetProcessInstanceList").Logger
	_, err := camunda.Connection.Delete("/process-instance/" + Id)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
