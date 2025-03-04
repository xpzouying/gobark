package gobark_test

import (
	"context"
	"log"

	"github.com/zy/gobark"
)

func Example() {
	// Create a new Bark client
	client, err := gobark.NewClient("https://api.day.app", "YOUR_BARK_KEY")
	if err != nil {
		log.Fatal(err)
	}

	// Send a simple notification with just body
	err = client.Send(context.Background(), "Hello from Bark!")
	if err != nil {
		log.Fatal(err)
	}

	// Send a notification with title and body
	err = client.Send(context.Background(), "This is the message body",
		gobark.WithTitle("Custom Title"))
	if err != nil {
		log.Fatal(err)
	}

	// Send a notification with all options
	err = client.Send(context.Background(), "Important message!",
		gobark.WithTitle("Meeting Reminder"),
		gobark.WithSubtitle("Team Standup"),
		gobark.WithIcon("https://example.com/icon.png"),
		gobark.WithSound("bell"),
		gobark.WithTimeSensitive(),
		gobark.WithCriticalNotify(),
	)
	if err != nil {
		log.Fatal(err)
	}
}
