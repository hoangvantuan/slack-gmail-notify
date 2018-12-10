package worker

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/nlopes/slack"
)

type slWorker interface {
	post(m *message, to string) error
	posts(m []*message, to string) error
}

type slWorkerImpl struct {
	client *slack.Client
}

func newSlWorker(c *slack.Client) slWorker {
	return &slWorkerImpl{c}
}

func (n *slWorkerImpl) post(m *message, to string) error {
	_, err := n.client.UploadFile(slack.FileUploadParameters{
		Filetype: "post",
		Channels: []string{to},
		Content:  fmt.Sprintf("### FROM: %s\n### CC: %s\n\n%s", m.From, m.CC, m.Body),
		Filename: m.Subject,
	})
	if err != nil {
		return errors.Wrap(err, "have error while post message")
	}

	return nil
}

func (n *slWorkerImpl) posts(ms []*message, to string) error {
	for _, m := range ms {
		err := n.post(m, to)
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}

	return nil
}
