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
		s.ProcessMessage(d.Body)
	}
}

func (s *RabbitMQService) ProcessMessage(msg []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(msg, &event); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
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
			Category:       courseData["category"].(string),
			ImageURL:       courseData["image_url"].(string),
			Duration:       int(courseData["duration"].(float64)),
			AvailableSeats: int(courseData["available_seats"].(float64)),
		}
		s.courseService.UpdateCourse(course)
	case "delete":
		courseID := courseData["id"].(string)
		log.Printf("Deleting course with ID: %s", courseID)
		if err := s.courseService.DeleteCourse(courseID); err != nil {
			log.Printf("Error deleting course from Solr: %v", err)
		} else {
			log.Printf("Successfully deleted course from Solr")
		}
	}
}
