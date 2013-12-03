package goblueboxapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBlockParams_Validates(t *testing.T) {
	p := BlockParams{}

	err := p.Validates()
	if err == nil {
		t.Error("BlockParams.Validates() expected to return an error with no params")
	}

	p.Product = "foobar"
	p.Template = "foobar"
	p.Password = "foobar"
	p.SshPublicKey = "foobar"
	err = p.Validates()
	if err == nil {
		t.Error("BlockParams.Validates() expected to return an error with both password and ssh public key set")
	}
}

func TestBlocksService_List(t *testing.T) {
	setup()
	defer teardown()

	output := `[
		{"id": "abcdef", "hostname": "abcdef.example.com", "ips":[{"address":"127.0.0.1"}, {"address": "::1"}], "status":"running"},
		{"id": "ghijkl", "hostname": "ghijkl.example.com", "ips":[{"address":"10.0.0.1"}], "status": "queued"}
	]`

	mux.HandleFunc("/api/blocks.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, output)
	})

	blocks, err := client.Blocks.List()

	if err != nil {
		t.Errorf("List() err expected to be nil, was %v", err)
	}

	want := []Block{
		Block{
			Id:       "abcdef",
			Hostname: "abcdef.example.com",
			Ips: []BlockIp{
				BlockIp{Address: "127.0.0.1"},
				BlockIp{Address: "::1"},
			},
			Status: "running",
		},
		Block{
			Id:       "ghijkl",
			Hostname: "ghijkl.example.com",
			Ips: []BlockIp{
				BlockIp{Address: "10.0.0.1"},
			},
			Status: "queued",
		},
	}

	if !reflect.DeepEqual(blocks, want) {
		t.Errorf("Blocks.List() returned %+v, want %+v", blocks, want)
	}
}

func TestBlocksService_Get(t *testing.T) {
	setup()
	defer teardown()

	output := `{"id": "abcdef", "hostname": "abcdef.example.com", "ips":[{"address":"127.0.0.1"}, {"address": "::1"}], "status":"running"}`

	mux.HandleFunc("/api/blocks/abcdef.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, output)
	})

	block, err := client.Blocks.Get("abcdef")

	if err != nil {
		t.Errorf("List() err expected to be nil, was %v", err)
	}

	want := &Block{
		Id:       "abcdef",
		Hostname: "abcdef.example.com",
		Ips: []BlockIp{
			BlockIp{Address: "127.0.0.1"},
			BlockIp{Address: "::1"},
		},
		Status: "running",
	}

	if !reflect.DeepEqual(block, want) {
		t.Errorf("Blocks.Get() returned %+v, want %+v", block, want)
	}
}

func TestBlocksService_Create(t *testing.T) {
	setup()
	defer teardown()

	output := `{"id": "ghijkl", "hostname": "ghijkl.example.com", "ips":[{"address":"10.0.0.1"}], "status": "queued"}`

	mux.HandleFunc("/api/blocks.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		if r.FormValue("product") != "the-product" || r.FormValue("template") != "the-template" || r.FormValue("password") != "the-password" {
			t.Error("Blocks.Create() expected to send params, but didn't")
		}

		fmt.Fprintf(w, output)
	})

	params := BlockParams{
		Product:  "the-product",
		Template: "the-template",
		Password: "the-password",
	}

	block, err := client.Blocks.Create(params)

	if err != nil {
		t.Errorf("Blocks.Create() returned error: %v", err)
	}

	want := &Block{
		Id:       "ghijkl",
		Hostname: "ghijkl.example.com",
		Ips: []BlockIp{
			BlockIp{Address: "10.0.0.1"},
		},
		Status: "queued",
	}

	if !reflect.DeepEqual(block, want) {
		t.Errorf("Blocks.Create() returned %+v, want %+v", block, want)
	}
}

func TestBlocksService_Destroy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/blocks/abcdef.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	err := client.Blocks.Destroy("abcdef")
	if err != nil {
		t.Errorf("Blocks.Destroy() returned error: %v", err)
	}
}
