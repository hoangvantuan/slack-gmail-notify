package worker

import (
	"encoding/base64"

	gm "google.golang.org/api/gmail/v1"
)

type ggWorker interface {
	fetchUnread() (*messages, error)
	read(m *message) error
	markLabel(m *message, labelID string) error
}

type ggWorkerImpl struct {
	srv   *gm.Service
	email string
}

func newGGWorker(srv *gm.Service, email string) ggWorker {
	return &ggWorkerImpl{srv, email}
}

func (g *ggWorkerImpl) fetchUnread() (*messages, error) {
	mrs, err := g.srv.Users.Messages.List("me").LabelIds("UNREAD").Q("NOT label:SLGMAILS").Do()
	if err != nil {
		return nil, err
	}

	return g.parseMessage(mrs)
}

func (g *ggWorkerImpl) read(m *message) error {
	_, err := g.srv.Users.Messages.Modify("me", m.ID, &gm.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}).Do()
	if err != nil {
		return err
	}

	return nil
}
func (g *ggWorkerImpl) markLabel(m *message, labelID string) error {
	_, err := g.srv.Users.Messages.Modify("me", m.ID, &gm.ModifyMessageRequest{
		AddLabelIds: []string{labelID},
	}).Do()
	if err != nil {
		return err
	}

	return nil
}

func (g *ggWorkerImpl) parseMessage(mr *gm.ListMessagesResponse) (*messages, error) {

	ms := &messages{}
	for _, m := range mr.Messages {
		msg := &message{}
		ms.ids = append(ms.ids, m.Id)
		msg.ID = m.Id
		msg.To = g.email

		msd, err := g.srv.Users.Messages.Get("me", m.Id).Do()
		if err != nil {
			return nil, err
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
			var body string

			if len(msd.Payload.Parts) == 0 {
				body = msd.Payload.Body.Data
			}

			for _, p := range msd.Payload.Parts {
				if p.MimeType == "text/plain" {
					body = p.Body.Data
				}
			}

			c, err := base64.URLEncoding.DecodeString(body)
			if err != nil {
				return nil, err
			}

			msg.Body = string(c)
		default:
			msg.Body = msd.Snippet
		}

		ms.m = append(ms.m, msg)
	}

	return ms, nil
}
