package tsmon

import (
	"context"
	"time"

	"golang.org/x/exp/slog"
)

type NotificationService interface {
	Send(ctx context.Context, notification Notification)
}

type Notifier struct {
	logger              *slog.Logger
	notificationQueue   <-chan Notification
	notificationService NotificationService
	lastNotification    time.Time
	offlineCounts       map[string]int
}

type NotifierOption func(n *Notifier)

func NewNotifier(opts ...NotifierOption) *Notifier {
	noti := &Notifier{
		logger:        slog.Default(),
		offlineCounts: map[string]int{},
	}
	for _, o := range opts {
		o(noti)
	}
	return noti
}

func WithQueue(ch <-chan Notification) NotifierOption {
	return func(n *Notifier) {
		n.notificationQueue = ch
	}
}

func WithNotificationService(svc NotificationService) NotifierOption {
	return func(n *Notifier) {
		n.notificationService = svc
	}
}

func (n *Notifier) Start() {
	for {
		event := <-n.notificationQueue
		go n.handleNotification(event)
	}
}

func (n *Notifier) reset(deviceNodeID string) {
	delete(n.offlineCounts, deviceNodeID)
}

func (n *Notifier) handleNotification(notification Notification) {
	if notification.Device == nil {
		n.logger.Error("nil device on notification")
		return
	}

	dev := notification.Device

	if dev.Online {
		n.reset(dev.NodeID)
		return
	}

	n.logger.Info("device offline", "device", dev.Name)

	if _, ok := n.offlineCounts[notification.Device.NodeID]; !ok {
		n.offlineCounts[notification.Device.NodeID] = 0
	}

	n.offlineCounts[notification.Device.NodeID]++

	deviceOfflineCount := n.offlineCounts[notification.Device.NodeID]
	n.logger.Info("consecutive offlines", "count", deviceOfflineCount, "device", dev.Name)
	// only if 5 in a row
	if time.Since(n.lastNotification) > time.Hour && deviceOfflineCount == 5 {
		n.logger.Info("sending notification to service")
		n.lastNotification = notification.Timestamp
		// todo need confirmation notification was successful
		n.notificationService.Send(context.TODO(), notification)
	}
}
