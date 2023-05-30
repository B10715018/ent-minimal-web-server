package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/b10715018/ent-minimal-web-server/ent/enttest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/steinfletcher/apitest"
)

func TestEntMinimalWebServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EntMinimalWebServer Suite")
}

var _ = Describe("Test Main", func() {
	It("should run", func() {
		client := enttest.Open(GinkgoT(), "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		defer client.Close()

		// seed the database with our "Hello, world" post and user.
		err := Seed(context.Background(), client)
		Expect(err).To(BeNil())

		srv := NewServer(client)
		r := NewRouter(srv)

		apitest.New().
			Handler(r).
			Get("/").
			Expect(GinkgoT()).
			Assert(func(resp *http.Response, _ *http.Request) error {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				x := "Hello, World!"
				if !strings.Contains(string(body), x) {
					return fmt.Errorf("Got : %s \n Want : %s", string(body), x)
				}
				return nil
			}).
			Status(http.StatusOK).
			End()
	})
})
