package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/vloldik/tasque"
	"github.com/vloldik/tasquelite"
)

type MailerContextKey string

var MailerContextValueKey = MailerContextKey("mailer")

// Type for e-mail task
type EmailTaskData struct {
	ID        int `gorm:"primaryKey"`
	Recipient string
	Subject   string
	Body      string
}

type Sender struct {
	wg *sync.WaitGroup
}

func (sender Sender) Send(value string) {
	time.Sleep(time.Second)
	fmt.Printf("sending value %s\n", value)
	sender.wg.Done()
}

// Do sends the email. This is a placeholder for the actual email sending logic.
func Do(ctx context.Context, emailTaskData EmailTaskData) {
	mailer := ctx.Value(MailerContextValueKey)
	if sender, ok := mailer.(*Sender); ok {
		sender.Send(emailTaskData.Body)
	} else {
		panic(errors.New("invalid mailer"))
	}
}

func main() {
	// Create a new context with a background.
	ctx := context.Background()

	// Create a new WaitGroup to wait for all tasks to complete.
	wg := &sync.WaitGroup{}

	// Create a new task queue with a capacity of 5.
	queue := tasque.NewTasksQueue(Do, 5)

	// Create a new sender with the WaitGroup.
	sender := Sender{wg: wg}

	// Add the sender to the context.
	ctx = context.WithValue(ctx, MailerContextValueKey, &sender)

	// Create a new GORM task storage manager for EmailTask.
	storage, err := tasquelite.NewGormTaskStorageManager("test.db", &EmailTaskData{}, 2)
	queue.SetTaskStoreManager(storage)

	if err != nil {
		// Handle the error if storage creation fails.
		panic(err)
	}

	// Define a target email task.
	target := EmailTaskData{
		Recipient: "foo@bar.com",
		Subject:   "<div>Hello!</div>",
		Body:      "Some important Subj",
	}

	// Fill the queue with 5 tasks.
	for i := 0; i < 5; i++ {
		// Send the task to the queue.
		// ID will be 0
		queue.SendToQueue(ctx, target)
		// Increment the WaitGroup counter for each task.
		sender.wg.Add(1)
	}

	// Try to send 10 tasks to the queue, which will be stored in memory.
	for i := 0; i < 10; i++ {
		// Send the task to the queue.
		// ID from storage
		queue.SendToQueue(ctx, target)
		// Increment the WaitGroup counter for each task.
		sender.wg.Add(1)
	}
	queue.StartQueue(ctx)
	// Wait for all tasks to complete.
	wg.Wait()
}
