package repository

import (
	"testing"

	"github.com/mdshun/slack-gmail-notify/infra"
)

func init() {
	infra.Setup()
}

func TestAdd(t *testing.T) {
	user := &User{
		UserID: "12345",
		TeamID: "122345",
	}

	repo := NewUserRepository(infra.RDB)

	user, err := repo.Add(user)

	if err != nil {
		t.Errorf("[Warn] %s", err)
	}

	repo.DeleteByID(user.ID)
}

func TestDelete(t *testing.T) {
	user := &User{
		UserID: "12345",
		TeamID: "122345",
	}

	repo := NewUserRepository(infra.RDB)

	user, err := repo.Add(user)

	if err != nil {
		t.Errorf("[Warn] %s", err)
		return
	}

	err = repo.DeleteByID(user.ID)

	if err != nil {
		t.Errorf("[Warn] %s", err)
	}
}

func TestFindByID(t *testing.T) {
	user := &User{
		UserID: "12345",
		TeamID: "122345",
	}

	repo := NewUserRepository(infra.RDB)

	user, err := repo.Add(user)

	if err != nil {
		t.Errorf("[Warn] %s", err)
		return
	}

	result, err := repo.FindByID(user.ID)

	if err != nil && user != result {
		t.Errorf("[Warn] %s", err)
	}

	repo.DeleteByID(user.ID)
}
