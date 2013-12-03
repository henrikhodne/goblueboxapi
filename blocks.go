package goblueboxapi

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// BlocksService exposes the API endpoints to interact with blocks.
type BlocksService struct {
	client *Client
}

// A Block is an on demand virtual computing resource.
type Block struct {
	Id       string
	Hostname string
	Ips      []BlockIp
	Status   string
}

type BlockIp struct {
	Address string
}

type BlockParams struct {
	Product      string
	Template     string
	Password     string
	SshPublicKey string
	Hostname     string
	Username     string
	Location     string
}

func (p BlockParams) Validates() error {
	if p.Product == "" {
		return errors.New(`Must specify "Product"`)
	}
	if p.Template == "" {
		return errors.New(`Must specify "Template"`)
	}
	if p.Password != "" && p.SshPublicKey != "" {
		return errors.New(`Only one of "Password" and "SshPublicKey" may be specified`)
	}
	if p.Password == "" && p.SshPublicKey == "" {
		return errors.New(`One of "Password" and "SshPublicKey" must be specified`)
	}

	return nil
}

func (p BlockParams) ToValues() url.Values {
	v := url.Values{}
	if p.Product != "" {
		v.Set("product", p.Product)
	}
	if p.Template != "" {
		v.Set("template", p.Template)
	}
	if p.Password != "" {
		v.Set("password", p.Password)
	}
	if p.SshPublicKey != "" {
		v.Set("ssh_public_key", p.SshPublicKey)
	}
	if p.Hostname != "" {
		v.Set("hostname", p.Hostname)
	}
	if p.Username != "" {
		v.Set("username", p.Username)
	}
	if p.Location != "" {
		v.Set("location", p.Location)
	}

	return v
}

func (s *BlocksService) List() ([]Block, error) {
	req, err := s.client.NewRequest("GET", "/blocks", nil)
	if err != nil {
		return nil, err
	}

	var blocks []Block
	err = s.client.Do(req, &blocks)

	return blocks, err
}

func (s *BlocksService) Get(uuid string) (*Block, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("/blocks/%s", uuid), nil)
	if err != nil {
		return nil, err
	}

	block := new(Block)
	err = s.client.Do(req, block)

	return block, err
}

func (s *BlocksService) Create(p BlockParams) (*Block, error) {
	if err := p.Validates(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", "/blocks", strings.NewReader(p.ToValues().Encode()))
	if err != nil {
		return nil, err
	}

	block := new(Block)
	err = s.client.Do(req, block)

	return block, err
}

func (s *BlocksService) Destroy(uuid string) error {
	req, err := s.client.NewRequest("DELETE", fmt.Sprintf("/blocks/%s", uuid), nil)
	if err != nil {
		return err
	}

	return s.client.Do(req, nil)
}
