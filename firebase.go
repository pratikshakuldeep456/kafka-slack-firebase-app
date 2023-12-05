package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func sendFirebaseNotification(message string) {
	opt := option.WithCredentialsFile(firebaseConfig)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("Error initializing Firebase app: %v", err)
		return
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Printf("Error creating Firebase messaging client: %v", err)
		return
	}

	msg := &messaging.Message{
		Data: map[string]string{
			"body": message,
		},
		Token: "firebase_device_token",
	}

	response, err := client.Send(context.Background(), msg)
	if err != nil {
		log.Printf("Error sending Firebase message: %v", err)
		return
	}

	fmt.Printf("Successfully sent message to Firebase: %v\n", response)
}
