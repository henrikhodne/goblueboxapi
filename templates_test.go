package goblueboxapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestTemplatesService_List(t *testing.T) {
	setup()
	defer teardown()

	now := time.Now()
	jsonNow, _ := now.MarshalText()

	output := `[
		{"id": "abcdef", "description": "foo bar baz", "public": true, "created": "%s"},
		{"id": "abcdefg", "description": "foo bar baz boom", "public": false, "created": "%s"}
	]`

	mux.HandleFunc("/api/block_templates.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, fmt.Sprintf(output, jsonNow, jsonNow))
	})

	templates, err := client.Templates.List()

	if err != nil {
		t.Errorf("List() err expected to be nil, was %v", err)
	}

	want := []Template{
		Template{
			ID:          "abcdef",
			Description: "foo bar baz",
			Public:      true,
			Created:     now,
		},
		Template{
			ID:          "abcdefg",
			Description: "foo bar baz boom",
			Public:      false,
			Created:     now,
		},
	}

	if !reflect.DeepEqual(templates, want) {
		t.Errorf("Templates.List() returned %+v, want %+v", templates, want)
	}
}

func TestTemplatesService_Get(t *testing.T) {
	setup()
	defer teardown()

	now := time.Now()
	jsonNow, _ := now.MarshalText()

	output := `{"id": "abcdef", "description": "foo bar baz", "public": true, "created": "%s"}`

	mux.HandleFunc("/api/block_templates/abcdef.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, fmt.Sprintf(output, jsonNow))
	})

	block, err := client.Templates.Get("abcdef")

	if err != nil {
		t.Errorf("List() err expected to be nil, was %v", err)
	}

	want := &Template{
		ID:          "abcdef",
		Description: "foo bar baz",
		Public:      true,
		Created:     now,
	}

	if !reflect.DeepEqual(block, want) {
		t.Errorf("Templates.Get() returned %+v, want %+v", block, want)
	}
}

func TestTemplatesService_Create(t *testing.T) {
	setup()
	defer teardown()

	output := `{"id": "abcdef", "description": "foo bar baz", "public": true, "created": "%s"}`

	mux.HandleFunc("/api/block_templates.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// if r.FormValue("product") != "the-product" ||
		// 	r.FormValue("template") != "the-template" ||
		// 	r.FormValue("password") != "the-password" ||
		// 	r.FormValue("ipv6_only") != "true" {
		// 	t.Error("Blocks.Create() expected to send params, but didn't")
		// }

		fmt.Fprintf(w, output)
	})

	// params := BlockParams{
	// 	Product:  "the-product",
	// 	Template: "the-template",
	// 	Password: "the-password",
	// 	IPv6Only: true,
	// }

	template, err := client.Templates.Create("abcdefg")

	if err != nil {
		t.Errorf("Templates.Create() returned error: %v", err)
	}

	want := &TemplateCreationStatus{
		Status: "abcdef",
		Text:   "queued",
		Error:  0,
	}

	if !reflect.DeepEqual(template, want) {
		t.Errorf("Templates.Create() returned %+v, want %+v", template, want)
	}
}

func TestTemplatesService_Destroy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/block_templates/abcdef.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	err := client.Templates.Destroy("abcdef")
	if err != nil {
		t.Errorf("Blocks.Destroy() returned error: %v", err)
	}
}
