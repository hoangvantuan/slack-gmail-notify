package util

import (
	"testing"
)

var key = "DuI9ucDupmPpU7C@jzdXv5AZw9^mlX*@"

func TestEncrypt(t *testing.T) {
	data := map[string]string{
		"value":  "testtest",
		"expect": "fc39931030493ae4a5d249306afb04bbdcc9fb69b06b30fb",
	}

	result, err := Encrypt(data["value"], key)

	if err != nil {
		t.Fatalf("[Error] %s", err)
	}

	if data["expect"] != string(result) {
		t.Errorf("[Error] expect %s receive %s", data["expect"], result)
	} else {
		t.Log(result)
	}
}

func TestDecrypt(t *testing.T) {
	data := map[string]string{
		"value":  "fc39931030493ae4a5d249306afb04bbdcc9fb69b06b30fb",
		"expect": "testtest",
	}

	result, err := Decrypt(data["value"], key)

	if err != nil {
		t.Errorf("[Error] %s", err)
	}

	if data["expect"] != string(result) {
		t.Errorf("[Error] expect %s receive %s", data["expect"], string(result))
	} else {
		t.Log(string(result))
	}
}
