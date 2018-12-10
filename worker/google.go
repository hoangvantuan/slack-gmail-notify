package worker

import (
	"encoding/base64"

	"github.com/pkg/errors"
	gm "google.golang.org/api/gmail/v1"
)

type ggWorker interface {
	fetchUnread() (*messages, error)
	read(m *messages) error
}

type ggWorkerImpl struct {
	srv *gm.Service
}

func newGGWorker(srv *gm.Service) ggWorker {
	return &ggWorkerImpl{srv}
}

func (g *ggWorkerImpl) fetchUnread() (*messages, error) {
	mrs, err := g.srv.Users.Messages.List("me").LabelIds("UNREAD").Do()
	if err != nil {
		return nil, errors.Wrap(err, "can not fetch list message")
	}

	return g.parseMessage(mrs), nil
}

func (g *ggWorkerImpl) read(m *messages) error {
	if m == nil || len(m.ids) == 0 {
		return nil
	}

	err := g.srv.Users.Messages.BatchModify("me", &gm.BatchModifyMessagesRequest{
		Ids:            m.ids,
		RemoveLabelIds: []string{"UNREAD"},
	}).Do()
	if err != nil {
		return errors.Wrap(err, "can not remove UNREAD label")
	}

	return nil
}

func (g *ggWorkerImpl) parseMessage(mr *gm.ListMessagesResponse) *messages {
	ms := &messages{}
	for _, m := range mr.Messages {
		msg := &message{}
		ms.ids = append(ms.ids, m.Id)

		msd, err := g.srv.Users.Messages.Get("me", m.Id).Do()
		if err != nil {
			return &messages{}
		}

		// parse header
		for _, h := range msd.Payload.Headers {
			if h.Name == "From" {
				msg.From = h.Value
			}

			if h.Name == "Subject" {
				msg.Subject = h.Value
			}

			if h.Name == "Cc" {
				msg.CC = h.Value
			}
		}
		// parse body
		switch {
		case msd.Payload.MimeType == "text/plain" || msd.Payload.MimeType == "multipart/alternative":
			for _, p := range msd.Payload.Parts {
				if p.MimeType == "text/plain" {
					c, err := base64.URLEncoding.DecodeString(p.Body.Data)
					if err != nil {
						return &messages{}
					}

					msg.Body = string(c)
				}
			}
		default:
			msg.Body = msd.Snippet
		}

		ms.m = append(ms.m, msg)
	}

	return ms
}
