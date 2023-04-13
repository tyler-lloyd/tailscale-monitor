package main

import (
	"flag"
	"strings"
	"tsmon"
	"tsmon/twilio"
)

func main() {
	var devices string
	var notificationBackendConfig string
	flag.StringVar(&devices, "devices", "", "list of devices to watch")
	flag.StringVar(&notificationBackendConfig, "config", "", "config path for notification backend service")
	flag.Parse()

	deviceList := strings.Split(devices, ",")
	notifyCh := make(chan tsmon.Notification)
	watcher := tsmon.NewWatchProcess(tsmon.NewLocalClient(),
		tsmon.WithNodesToWatch(deviceList),
		tsmon.WithNotificationChan(notifyCh),
	)

	notifier := tsmon.NewNotifier(
		tsmon.WithQueue(notifyCh),
		tsmon.WithNotificationService(
			twilio.NewTwilioService(notificationBackendConfig),
		),
	)
	go notifier.Start()

	watcher.Run()
}
