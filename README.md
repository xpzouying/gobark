# gobark

A Go SDK for [Bark](https://github.com/Finb/Bark) - a simple and secure iOS push notification service.

## Features

- Simple and idiomatic Go API
- Support for all major Bark notification features:
  - Custom titles and subtitles
  - Custom icons (iOS 15+)
  - Custom notification sounds
  - Time-sensitive notifications
  - Critical alerts
- Context support for cancellation and timeouts
- Configurable base URL for self-hosted Bark servers

## Installation

```bash
go get github.com/xpzouying/gobark
```

## Usage

```go
package main

import (
    "context"
    "log"

    "github.com/xpzouying/gobark"
)

func main() {
    // Create a new client
    client, err := gobark.NewClient("https://api.day.app", "YOUR_BARK_KEY")
    if err != nil {
        log.Fatal(err)
    }

    // Send a simple notification
    err = client.Send(context.Background(), "Hello from Bark!")
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
```

## Available Options

- `WithTitle(title string)`: Set notification title
- `WithSubtitle(subtitle string)`: Set notification subtitle
- `WithIcon(iconURL string)`: Set notification icon (iOS 15+ only)
- `WithSound(sound string)`: Set notification sound
- `WithTimeSensitive()`: Mark notification as time-sensitive
- `WithCriticalNotify()`: Mark notification as critical alert

## License

MIT License
