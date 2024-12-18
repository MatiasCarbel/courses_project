package services

import (
	"courses-api/models"
	"encoding/json"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	PublishCourseUpdate(course *models.Course, action string) error
	PublishCourseDelete(courseID interface{}) error
}

type RabbitMQService struct {
	uri string
}

func NewRabbitMQService() MessageQueue {
	uri := os.Getenv("RABBITMQ_URI")
	if uri == "" {
		uri = "amqp://guest:guest@rabbitmq:5672/"
	}
	return &RabbitMQService{uri: uri}
}

func (r *RabbitMQService) PublishCourseUpdate(course *models.Course, action string) error {
	conn, err := amqp091.Dial(r.uri)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"course_updates",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	event := map[string]interface{}{
		"action": action,
		"course": course,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (r *RabbitMQService) PublishCourseDelete(courseID interface{}) error {
	conn, err := amqp091.Dial(r.uri)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"course_updates",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	event := map[string]interface{}{
		"action": "delete",
		"course": map[string]interface{}{
			"id": courseID,
		},
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
} 