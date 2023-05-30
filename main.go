package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/b10715018/ent-minimal-web-server/ent"
	"github.com/b10715018/ent-minimal-web-server/ent/post"
	"github.com/b10715018/ent-minimal-web-server/ent/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("HOST_NAME")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("DB_NAME")
	if host == "" || port == "" || user == "" || dbName == "" {
		log.Fatalf("need to give host, port, user and db name")
	}
	client, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, user, dbName))
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	if !client.Post.Query().ExistX(ctx) {
		if err := Seed(ctx, client); err != nil {
			log.Fatalf("failed seeding the database: %v", err)
		}
	}
	srv := NewServer(client)
	r := NewRouter(srv)
	log.Fatal(http.ListenAndServe(":8080", r))

}

func NewRouter(srv *server) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", srv.index)
	r.Post("/add", srv.add)
	return r
}

var (
	//go:embed templates/*
	resources embed.FS
	tmpl      = template.Must(template.ParseFS(resources, "templates/*"))
)

type server struct {
	client *ent.Client
}

func NewServer(client *ent.Client) *server {
	return &server{client: client}
}

// index serves the blog home page
func (s *server) index(w http.ResponseWriter, r *http.Request) {
	posts, err := s.client.Post.
		Query().
		WithAuthor().
		Order(ent.Desc(post.FieldCreatedAt)).
		All(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) add(w http.ResponseWriter, r *http.Request) {
	author, err := s.client.User.Query().Only(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := s.client.Post.Create().
		SetTitle(r.FormValue("title")).
		SetBody(r.FormValue("body")).
		SetAuthor(author).
		Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func Seed(ctx context.Context, client *ent.Client) error {
	// Check if the user "brandon" already exists.
	r, err := client.User.Query().
		Where(
			user.Name("brandon"),
		).
		Only(ctx)
	switch {
	// If not, create the user.
	case ent.IsNotFound(err):
		r, err = client.User.Create().
			SetName("brandon").
			SetEmail("test@mail.com").
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed creating user: %v", err)
		}
	case err != nil:
		return fmt.Errorf("failed querying user: %v", err)
	}
	// Finally, create a "Hello, world" blogpost.
	return client.Post.Create().
		SetTitle("Hello, World!").
		SetBody("This is my first post").
		SetAuthor(r).
		Exec(ctx)
}
