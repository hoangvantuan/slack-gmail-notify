package worker

import (
	"fmt"

	"github.com/nlopes/slack"
)

type slWorker interface {
	post(gg ggWorker, m *message, to string) error
	posts(gg ggWorker, m []*message, to string) error
}

type slWorkerImpl struct {
	client *slack.Client
}

func newSlWorker(c *slack.Client) slWorker {
	return &slWorkerImpl{c}
}

func (n *slWorkerImpl) post(gg ggWorker, m *message, to string) error {
	_, err := n.client.UploadFile(slack.FileUploadParameters{
		Filetype: "post",
		Channels: []string{to},
		Content:  fmt.Sprintf("### FROM: %s\n\n%s", m.From, m.Body),
		Filename: m.Subject,
	})
	if err != nil {
		return err
	}

	err = gg.read(m)
	if err != nil {
		return err
	}

	return nil
}

func (n *slWorkerImpl) posts(gg ggWorker, ms []*message, to string) error {
	for _, m := range ms {
		err := n.post(gg, m, to)
		if err != nil {
			return err
		}
		//time.Sleep(time.Second)
	}

	return nil
}
