package main

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"fmt"
	"errors"
	"encoding/json"
	"os/signal"
	"os"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"log"
)

const topicName = "default-topic"
const BrokerAddr = "0.0.0.0:51015"
const processFileBpmn = "src/simpleProcess.bpmn"
const processFileYaml = "src/simpleProcess.yml"
const processId = "simpleProcess"
const taskA = "task-a"

var errClientStartFailed = errors.New("cannot start client")
var errWorkflowDeploymentFailed = errors.New("creating new workflow deployment failed")

func main() {

	zbClient := createNewClient()

	createDefaultTopicIfNotExists(zbClient)

	outputTopologie(zbClient)

	deployProcessBpmn(zbClient)
	//deployProcessYaml(zbClient) <- is working but than you don't see a BPMN in Monitor

	startProcess(zbClient)

	subscription := createSubscriptionForTaskA(zbClient)

	startGoRoutineToCloseSubscriptionOnExit(subscription)

	subscription.Start()
}

func createNewClient() (*zbc.Client) {
	fmt.Println("Create new zeebe client")

	zbClient, err := zbc.NewClient(BrokerAddr)
	if err != nil {
		panic(errClientStartFailed)
	}

	return zbClient
}

func outputTopologie(zbClient *zbc.Client) {
	fmt.Println("Load broker topologie")

	topology, err := zbClient.RefreshTopology()
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(topology, "", "    ")
	fmt.Println("Topologie: ", string(b))
}

func createDefaultTopicIfNotExists(zbClient *zbc.Client) {
	log.Printf("Create new topic '%s'", topicName)

	if topicExists(zbClient, topicName) {
		log.Println("Topic does already exist")
		return
	}

	topic, err := zbClient.CreateTopic(topicName, 1)
	if err != nil {
		log.Fatal("Could not create topic")
		panic(err)
	}

	log.Println("Created topic: ", topic)
}

func topicExists(zbClient *zbc.Client, topicName string) bool {
	topology, err := zbClient.RefreshTopology()
	if err != nil {
		log.Fatal("Error happens while loading topology")
		panic(err)
	}

	return topology.PartitionIDByTopicName[topicName] != nil
}

func deployProcessBpmn(zbClient *zbc.Client) {
	deployProcess(zbClient, zbcommon.BpmnXml, processFileBpmn)
}

func deployProcessYaml(zbClient *zbc.Client) {
	deployProcess(zbClient, zbcommon.YamlWorkflow, processFileYaml)
}

func deployProcess(zbClient *zbc.Client, resourceType, path string) {
	fmt.Printf("Deploy '%s' process '%s'\n", resourceType, processFileYaml)

	response, err := zbClient.CreateWorkflowFromFile(topicName, resourceType, path)
	if err != nil {
		panic(errWorkflowDeploymentFailed)
	}

	fmt.Println("Deployed Process Responce: ", response.String())
}

func startProcess(zbClient *zbc.Client) {
	fmt.Println("Start process ", processId)

	payload := make(map[string]interface{})
	payload["somePayload"] = "31243"
	payload["someOtherPayload"] = "lol"

	instance := zbc.NewWorkflowInstance(processId, -1, payload)
	msg, err := zbClient.CreateWorkflowInstance(topicName, instance)

	if err != nil {
		panic(err)
	}

	fmt.Println("Start Process responce: ", msg.String())
}

func createSubscriptionForTaskA(zbClient *zbc.Client) *zbsubscribe.TaskSubscription {
	fmt.Println("Open task subscription for Task A")

	subscription, err := zbClient.TaskSubscription(topicName, "lockOwner", taskA, 32, func(clientApi zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
		fmt.Println("Message of task A subscription: ", event.String())

		// complete task after processing
		response, _ := clientApi.CompleteTask(event)
		fmt.Println("Complete Task Responce: ", response)
	})

	if err != nil {
		panic("Unable to open subscription")
	}

	return subscription
}

func startGoRoutineToCloseSubscriptionOnExit(subscription *zbsubscribe.TaskSubscription) {
	osCh := make(chan os.Signal, 1)
	signal.Notify(osCh, os.Interrupt)
	go func() {
		fmt.Println("Create go routine which waits for app interrrupt")

		<-osCh
		fmt.Println("Closing subscription.")
		err := subscription.Close()
		if err != nil {
			fmt.Println("Failed to close subscription: ", err)
		} else {
			fmt.Println("Subscription closed.")
		}
		os.Exit(0)
	}()
}
