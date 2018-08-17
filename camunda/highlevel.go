package camunda

import (
	"encoding/json"

	_ "github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func (camunda *CamundaConnection) GetProcessInstanceList() (ProcessInstanceList, error) {
	var list ProcessInstanceList

	logger := log.WithField("func", "GetProcessInstanceList").Logger
	raw, err := camunda.GetRaw("/process-instance")
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
	_, err := camunda.DeleteRaw("/process-instance/" + Id)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
