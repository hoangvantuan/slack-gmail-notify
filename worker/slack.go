package worker

import (
	"fmt"
	"time"

	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/nlopes/slack"
)

type slWorker interface {
	post(gg ggWorker, m *message, to string, isRead bool, labelID string) error
	posts(gg ggWorker, m []*message, to string, isRead bool, labelID string) error
}

type slWorkerImpl struct {
	client *slack.Client
}

func newSlWorker(c *slack.Client) slWorker {
	return &slWorkerImpl{c}
}

func (n *slWorkerImpl) post(gg ggWorker, m *message, to string, isRead bool, labelID string) error {
	_, err := n.client.UploadFile(slack.FileUploadParameters{
		Filetype: "post",
		Channels: []string{to},
		Content:  fmt.Sprintf("### To: %s\n### From: %s\n\n%s", m.To, m.From, m.Body),
		Filename: m.Subject,
		Title:    m.Subject,
	})
	if err != nil {
		return err
	}

	err = gg.markLabel(m, labelID)
	if err != nil {
		return err
	}

	if isRead {
		err = gg.read(m)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *slWorkerImpl) posts(gg ggWorker, ms []*message, to string, isRead bool, labelID string) error {
	count := 1
	for _, m := range ms {
		infra.Debug(fmt.Sprintf("%d message was sent to %s", count, to))
		count = count + 1
		err := n.post(gg, m, to, isRead, labelID)
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}

	return nil
}
