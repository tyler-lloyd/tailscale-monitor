package tsmon

import (
	"context"
	"time"

	"golang.org/x/exp/slog"
)

type TailscaleClient interface {
	GetDevice(ctx context.Context, id string) (*TailnetDevice, error)
	GetDevices(ctx context.Context) ([]*TailnetDevice, error)
}

type TailscaleWatchProcess struct {
	client     TailscaleClient
	nodes      []string
	notifyChan chan Notification
	logger     *slog.Logger
}

type TailscaleWatchOption func(*TailscaleWatchProcess)

func NewWatchProcess(c TailscaleClient, options ...TailscaleWatchOption) *TailscaleWatchProcess {
	proc := &TailscaleWatchProcess{
		client: c,
		logger: slog.Default(),
	}

	for _, opt := range options {
		opt(proc)
	}

	return proc
}

func WithNodesToWatch(ids []string) TailscaleWatchOption {
	return func(twp *TailscaleWatchProcess) {
		if twp.nodes == nil {
			twp.nodes = make([]string, 0)
		}
		twp.nodes = append(twp.nodes, ids...)
	}
}

func WithNotificationChan(ch chan Notification) TailscaleWatchOption {
	return func(twp *TailscaleWatchProcess) {
		twp.notifyChan = ch
	}
}

func WithLogger(l *slog.Logger) TailscaleWatchOption {
	return func(twp *TailscaleWatchProcess) {
		twp.logger = l
	}
}

func (t *TailscaleWatchProcess) Run() {
	ctx := context.Background()
	for {
		t.Poll(ctx)
		time.Sleep(10 * time.Second)
	}
}

func (t *TailscaleWatchProcess) Poll(ctx context.Context) {
	devs, err := t.client.GetDevices(ctx)
	if err != nil {
		t.logger.ErrorCtx(ctx, "failed to get devices", "error", err)
		return
	}

	deviceMap := make(map[string]*TailnetDevice)
	for i, d := range devs {
		deviceMap[d.ID] = devs[i]
	}

	for _, id := range t.nodes {
		d, ok := deviceMap[id]
		if !ok {
			t.logger.WarnCtx(ctx, "device not found", "device", id)
			continue
		}
		if !d.Online {
			t.logger.InfoCtx(ctx, "offline. notifying.", "device", d.Name, "lastSeen", d.LastSeen)
			go t.notifyOffline(d)
			continue
		}
		t.logger.InfoCtx(ctx, "online", "device", d.Name)
	}
}

func (t *TailscaleWatchProcess) notifyOffline(dev *TailnetDevice) {
	t.notifyChan <- Notification{dev, time.Now()}
}
