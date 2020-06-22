package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"todo-service/model"
)

type Repository interface {
	FindAll() interface{}
	Find(id string) interface{}
	Create(i interface{}) interface{}
	Update(i interface{}) interface{}
	Delete(id string)
}

type TodoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(collection *mongo.Collection) Repository {
	return TodoRepository{collection: collection}
}

func (r TodoRepository) FindAll() interface{} {
	var todos = make([]model.Todo, 0)
	found, _ := r.collection.Find(context.Background(), bson.D{})
	_ = found.All(context.Background(), &todos)
	return todos
}

func (r TodoRepository) Find(id string) interface{} {
	var todo model.Todo
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	one := r.collection.FindOne(context.Background(), filter)
	_ = one.Decode(&todo)
	return todo
}

func (r TodoRepository) Create(i interface{}) interface{} {
	id, _ := r.collection.InsertOne(context.Background(), i)
	objectId, _ := id.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": objectId}
	one := r.collection.FindOne(context.Background(), filter)
	var todo model.Todo
	_ = one.Decode(&todo)
	return todo
}

func (r TodoRepository) Update(i interface{}) interface{} {
	todo := i.(model.Todo)
	id, _ := primitive.ObjectIDFromHex(todo.ID)
	filter := bson.M{"_id": id}
	_, err := r.collection.UpdateOne(context.Background(), filter, bson.D{
		{"$set", bson.D{{"title", todo.Title}, {"completed", todo.Completed}}},
	})

	if err != nil {
		return err
	}

	one := r.collection.FindOne(context.Background(), filter)
	_ = one.Decode(&todo)
	return todo
}

func (r TodoRepository) Delete(id string) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	_, _ = r.collection.DeleteOne(context.Background(), filter)
}
