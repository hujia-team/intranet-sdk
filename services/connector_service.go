// Package services provides business logic for the MiniEye Intranet API.
package services

import (
	"encoding/json"

	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
	"github.com/hujia-team/intranet-sdk/utils"
)

type KafkaMessage struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

// ConnectorService defines the connector service interface.
type ConnectorService interface {
	// SendKafkaMessage sends a message to Kafka.
	SendKafkaMessage(topic string, message any) (models.BaseMsgResp, error)
}

// connectorService implements the ConnectorService interface.
type connectorService struct {
	httpClient *client.HTTPClient
}

// NewConnectorService creates a new connector service.
func NewConnectorService(httpClient *client.HTTPClient) ConnectorService {
	return &connectorService{
		httpClient: httpClient,
	}
}

// SendKafkaMessage implements the ConnectorService.SendKafkaMessage method.
func (s *connectorService) SendKafkaMessage(topic string, message any) (models.BaseMsgResp, error) {
	// 使用嵌套结构体直接解析响应
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	jsonStr, err := json.Marshal(message)
	if err != nil {
		utils.Error("JSON序列化失败: %v", err)
		return models.BaseMsgResp{}, utils.NewInternalError("failed to marshal message", err)
	}
	utils.Debug("Sending message to Kafka topic: %s", topic)
	err = s.httpClient.Post("/connector/kafka/send-topic-message", KafkaMessage{
		Topic:   topic,
		Message: string(jsonStr),
	}, &response)
	if err != nil {
		utils.Error("Failed to send message to Kafka: %v", err)
		return models.BaseMsgResp{}, utils.NewAPIError("failed to send message to Kafka", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return models.BaseMsgResp{
			Code: response.Code,
			Msg:  response.Msg,
		}, utils.NewAPIError(response.Msg, nil)
	}

	utils.Debug("Sent message to Kafka topic: %s successfully", topic)
	return models.BaseMsgResp{
		Code: response.Code,
		Msg:  response.Msg,
	}, nil
}
