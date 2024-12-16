package server

import (
	"context"
	"db-service/internal/models"
	producer "db-service/pkg/kafka"
	pb "db-service/pkg/proto"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Server struct {
	pb.UnimplementedDBServiceServer
	DB *gorm.DB
	RDB *redis.Client
}


func (s *Server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	var user models.User

	cacheKey := fmt.Sprintf("user:%d", req.Id)
	cachedUser, err := s.RDB.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		if err := s.DB.First(&user, req.Id).Error; err != nil {
			return nil, fmt.Errorf("user not found: %v", err)
		}

		userJSON, err := json.Marshal(user)
		if err != nil {
			log.Printf("Error marshaling user: %v", err)
			return nil, fmt.Errorf("error marshaling user")
		}
		s.RDB.Set(ctx, cacheKey, userJSON, 10*time.Minute)

	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch from Redis: %v", err)
	} else {
		log.Println("User found in cache")
		err := json.Unmarshal([]byte(cachedUser), &user)
		if err != nil {
			log.Printf("Error unmarshaling user: %v", err)
			return nil, fmt.Errorf("error unmarshaling user")
		}
	}

	return &pb.UserResponse{
		Id:    int32(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user := models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	producer, err := producer.CreateKafkaProducer()
	if err != nil {
		log.Printf("Failed to create Kafka producer: %v", err)
		return nil, fmt.Errorf("failed to create Kafka producer")
	}
	defer producer.Close()

	message := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"time":  time.Now().Format(time.RFC3339),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal Kafka message: %v", err)
		return nil, fmt.Errorf("failed to marshal Kafka message")
	}

	topic := "user.created"
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          messageBytes,
	}, nil)

	if err != nil {
		log.Printf("Failed to produce Kafka message: %v", err)
		return nil, fmt.Errorf("failed to produce Kafka message")
	}

	log.Println("Kafka message sent successfully")

	return &pb.UserResponse{
		Id:    int32(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}