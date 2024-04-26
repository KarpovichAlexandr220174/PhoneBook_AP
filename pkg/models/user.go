package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type User struct {
	Username  string    `bson:"username"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"createdAt"`
}

func SaveUser(user User) error {
	// Ваша логика сохранения пользователя в базу данных NoSQL
	// Здесь можно использовать MongoDB или любую другую NoSQL базу данных

	// Пример для MongoDB:
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("usersdb").Collection("users")
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(emailOrUsername string) (User, error) {
	// Подключаем к базе данных
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return User{}, err
	}
	defer client.Disconnect(context.Background())

	// Получаем коллекцию пользователей
	collection := client.Database("usersdb").Collection("users")

	// Ищем пользователя по почте или нику
	filter := bson.M{"$or": []bson.M{{"email": emailOrUsername}, {"username": emailOrUsername}}}
	var user User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func UpdatePasswordByEmail(email, newPassword string) error {
	// Подключаемся к MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	// Получаем коллекцию пользователей
	collection := client.Database("usersdb").Collection("users")

	// Подготавливаем фильтр для поиска пользователя по email
	filter := bson.M{"email": email}

	// Подготавливаем обновление для установки нового пароля
	update := bson.M{"$set": bson.M{"password": newPassword}}

	// Выполняем обновление документа пользователя
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
