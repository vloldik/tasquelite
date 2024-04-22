package main

import (
	"context"
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

func (sender Sender) Send(value int) {
	time.Sleep(time.Second)
	fmt.Printf("sending value %d\n", value)
	sender.wg.Done()
}

// Do sends the email. This is a placeholder for the actual email sending logic.

func main() {
	// Create a new context with a background.
	ctx := context.Background()

	// Create a new WaitGroup to wait for all tasks to complete.
	wg := &sync.WaitGroup{}

	storage, err := tasquelite.NewGormTaskStorageManager("test.db", &EmailTaskData{})

	// Create a new task queue with a capacity of 5.

	// Create a new sender with the WaitGroup.
	sender := Sender{wg: wg}

	Do := func(ctx context.Context, emailTaskData EmailTaskData) bool {
		sender.Send(emailTaskData.ID)
		return false
	}

	queue := tasque.CreateTaskQueue(Do, 5, storage, func(err error) bool { fmt.Println(err); return false }, time.Hour)

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
