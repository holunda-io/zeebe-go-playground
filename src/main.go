package main

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"fmt"
	"errors"
	"encoding/json"
)

const topicName = "default-topic"
const BrokerAddr = "127.0.0.1:51015"
const processFile = "src/simpleProcess.bpmn"
const processId = "simpleProcess"

var errClientStartFailed = errors.New("cannot start client")
var errWorkflowDeploymentFailed = errors.New("creating new workflow deployment failed")

func main() {

	zbClient, err := zbc.NewClient(BrokerAddr)
	if err != nil {
		panic(errClientStartFailed)
	}

	loadTopologie(zbClient)

	deployProcess(zbClient)

	startProcess(zbClient)
}

func loadTopologie(zbClient *zbc.Client) {

	topology, err := zbClient.Topology()
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(topology, "", "    ")
	fmt.Println(string(b))
}

func deployProcess(zbClient *zbc.Client) {
	response, err := zbClient.CreateWorkflowFromFile(topicName, zbc.BpmnXml, processFile)
	if err != nil {
		panic(errWorkflowDeploymentFailed)
	}

	fmt.Println(response.String())
}

func startProcess(zbClient *zbc.Client) {
	payload := make(map[string]interface{})
	payload["somePayload"] = "31243"
	payload["someOtherPayload"] = "lol"

	instance := zbc.NewWorkflowInstance(processId, -1, payload)
	msg, err := zbClient.CreateWorkflowInstance(topicName, instance)

	if err != nil {
		panic(err)
	}

	fmt.Println(msg.String())
}
