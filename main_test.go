package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/b10715018/ent-minimal-web-server/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// seed the database with our "Hello, world" post and user.
	err := Seed(context.Background(), client)
	require.NoError(t, err)

	srv := NewServer(client)
	r := NewRouter(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "Hello, World!")
}

func TestAdd(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	err := Seed(context.Background(), client)
	require.NoError(t, err)

	srv := NewServer(client)
	r := NewRouter(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Post the form.
	resp, err := ts.Client().PostForm(ts.URL+"/add", map[string][]string{
		"title": {"Testing, one, two."},
		"body":  {"This is a test"},
	})
	require.NoError(t, err)
	// We should be redirected to the index page and receive 200 OK.
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// The home page should contain our new post.
	require.Contains(t, string(body), "This is a test")
}
