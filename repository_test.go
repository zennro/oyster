package main

import (
	"io/ioutil"
	"testing"

	"github.com/proglottis/rwvfs"
)

func TestFormRepoPutGet(t *testing.T) {
	repo := setupFormRepo(t)

	if _, err := repo.Get(&FormRequest{Key: "test"}, []byte("password")); err != ErrNotFound {
		t.Error("Expected ErrNotFound, got", err)
	}

	writeform := &Form{
		FormRequest: FormRequest{Key: "test"},
		Fields: map[string]string{
			"username": "bob",
			"password": "password123",
		},
	}
	if err := repo.Put(writeform); err != nil {
		t.Fatal(err)
	}

	readform, err := repo.Get(&FormRequest{Key: "test"}, []byte("password"))
	if err != nil {
		t.Fatal(err)
	}
	if readform.Key != "test" {
		t.Errorf("Expected 'test', got %#v", readform.Key)
	}
	for field, value := range writeform.Fields {
		if readform.Fields[field] != value {
			t.Errorf("Expected %#v, got %#v", value, readform.Fields[field])
		}
	}
}

func TestFileRepoCreateOpen(t *testing.T) {
	repo := setupFileRepo(t)

	if _, err := repo.Open("test", []byte("password")); err != ErrNotFound {
		t.Error("Expected ErrNotFound, got", err)
	}

	clearwrite, err := repo.Create("test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = clearwrite.Write([]byte("password123\nThe best password"))
	if err != nil {
		t.Fatal(err)
	}
	clearwrite.Close()

	clearread, err := repo.Open("test", []byte("password"))
	if err != nil {
		t.Fatal(err)
	}

	text, err := ioutil.ReadAll(clearread)
	if err != nil {
		t.Fatal(err)
	}

	clearread.Close()

	if string(text) != "password123\nThe best password" {
		t.Error("Expected 'password123\\nThe best password', got", string(text))
	}

	line, err := repo.Line("test", []byte("password"))
	if err != nil {
		t.Fatal(err)
	}

	if line != "password123" {
		t.Error("Expected 'password123', got", line)
	}
}

func setupFormRepo(t testing.TB) *FormRepo {
	gpg := NewGpgRepo("gpghome")
	fs := NewCryptoFS(rwvfs.Map(map[string]string{}), gpg)
	if err := InitRepo(fs, []string{"test@example.com"}); err != nil {
		t.Fatal(err)
	}
	return NewFormRepo(fs)
}

func setupFileRepo(t testing.TB) *FileRepo {
	gpg := NewGpgRepo("gpghome")
	fs := NewCryptoFS(rwvfs.Map(map[string]string{}), gpg)
	if err := InitRepo(fs, []string{"test@example.com"}); err != nil {
		t.Fatal(err)
	}
	return NewFileRepo(fs)
}
