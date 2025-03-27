package db

import (
	"SSO/internal/config"
	"SSO/internal/domain/models"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gorm.io/gorm"
	"log/slog"
	"sync"
)

const WorkerCount = 5

func SynchronizationUserTable(log *slog.Logger, params *config.Config) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%d", params.KafkaParams.Host, params.KafkaParams.Port),
		"group.id":          params.KafkaParams.GroupID,
		"auto.offset.reset": params.KafkaParams.AutoOffsetReset,
	})
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	err = consumer.Subscribe(params.KafkaParams.Topic, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Consumer started, waiting for messages...")

	// Канал для передачи задач в воркеры
	jobs := make(chan models.User, 100) // Буферизованный канал, чтобы не блокировать чтение

	// Запускаем воркеров
	var wg sync.WaitGroup
	for i := 0; i < WorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for user := range jobs {
				CreateUser(GetDBConn(params.AppParams.DBSM), log, user) // Создание пользователя в БД
			}
		}()
	}

	// Основной цикл обработки сообщений
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Error("Error reading kafka message", slog.String("error", err.Error()))

			continue
		}

		var user models.User
		if err = json.Unmarshal(msg.Value, &user); err != nil {
			log.Error("Error unmarshalling user data", slog.String("error", err.Error()))

			continue
		}

		// Отправляем задачу в канал
		jobs <- user
	}

	// Закрываем канал и ждем завершения всех воркеров
	close(jobs)
	wg.Wait()
}

func CreateUser(db *gorm.DB, log *slog.Logger, user models.User) (userDB models.User, err error) {
	//logger.Debug.Println(user.ID)
	if err = db.Create(&user).Error; err != nil {
		log.Error("Error creating user", slog.String("error", err.Error()))

		return userDB, err
	}

	//logger.Debug.Println(user.ID)
	userDB = user
	return userDB, nil
}
