package api

import (
	"google.golang.org/protobuf/proto"
	"k8s.io/klog/v2"
)

func createJob(message []byte) []byte {
	req := &CreateJobRequest{}
	proto.Unmarshal(message, req)

	jobResponse := &CreateJobResponse{}
	jobResponse.KubernetesJobName = "bla bla" //todo replace it

	marshal, err := proto.Marshal(jobResponse)

	if err != nil {
		klog.Infof("Marshal error: %s", err)
	}
	return marshal
}

func processingFunction() func(message []byte, functionId uint8) []byte {
	return func(message []byte, functionId uint8) []byte {
		klog.Infof("Pocessing function %s", functionId)
		if functionId == 1 {
			return createJob(message)
		}

		return []byte("FUNCTION_NOT_FOUND") //todo добавить error code
	}
}
