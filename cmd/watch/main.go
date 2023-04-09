package main

import (
	"flag"
	"strings"
	"tswatcher"
)

func main() {
	var devices string
	flag.StringVar(&devices, "devices", "", "list of devices to watch")
	flag.Parse()

	deviceList := strings.Split(devices, ",")
	notifyCh := make(chan tswatcher.Notification)
	watcher := tswatcher.NewWatchProcess(tswatcher.NewLocalClient(),
		tswatcher.WithNodesToWatch(deviceList),
		tswatcher.WithNotificationChan(notifyCh),
	)

	notifier := tswatcher.NewNotifier(
		tswatcher.WithQueue(notifyCh),
	)
	go notifier.Start()

	watcher.Run()
}
