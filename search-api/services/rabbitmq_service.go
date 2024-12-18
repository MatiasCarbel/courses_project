package services

import (
	"encoding/json"
	"log"

	"search-api/domain"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	courseService *CourseService
}

func NewRabbitMQService(courseService *CourseService) *RabbitMQService {
	return &RabbitMQService{courseService: courseService}
}

func (s *RabbitMQService) ConsumeRabbitMQMessages(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"course_updates",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	for d := range msgs {
		var event map[string]interface{}
		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		action := event["action"].(string)
		courseData := event["course"].(map[string]interface{})

		switch action {
		case "upsert":
			course := domain.Course{
				ID:             courseData["id"].(string),
				Title:          courseData["title"].(string),
				Description:    courseData["description"].(string),
				Instructor:     courseData["instructor"].(string),
				Duration:       int(courseData["duration"].(float64)),
				AvailableSeats: int(courseData["available_seats"].(float64)),
			}
			s.courseService.UpdateCourse(course)
		}
	}
}
