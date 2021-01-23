package manager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mar-tina/postmanager/schema"
)

type CollOpts struct {
	Schema      string `json:"schema"`
	Description string `json:"description"`
	Mode        string `json:"mode"`
}

func DefaultCollOpts() *CollOpts {
	return &CollOpts{
		Description: "no descritpion provided",
		Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		Mode:        "json",
	}
}

type Route struct {
	schema.Request
	name string
}

type Coll struct {
	collection   *schema.Collection
	routes       map[string]*Route
	currentroute string
	mode         string
}

func NewCollection(name string, opts *CollOpts) *Coll {
	return &Coll{
		collection: &schema.Collection{
			Info: schema.Info{
				Schema:      opts.Schema,
				Description: opts.Description,
			},
		},
		routes: make(map[string]*Route),
		mode:   opts.Mode,
	}
}

func (c *Coll) Route(name string) *Coll {
	c.routes[name] = &Route{
		name: name,
	}

	c.currentroute = name
	return c
}

func (c *Coll) Header(key, value string) *Coll {
	route := c.routes[c.currentroute]
	header := schema.Header{
		Key:   key,
		Value: value,
	}
	route.Header = append(route.Header, header)

	c.routes[c.currentroute] = route
	return c
}

func (c *Coll) HeaderWithEnv(key, value string) *Coll {
	route := c.routes[c.currentroute]
	value = fmt.Sprintf("{{%s}}", value)
	header := schema.Header{
		Key:   key,
		Value: value,
	}
	route.Header = append(route.Header, header)

	c.routes[c.currentroute] = route
	return c
}

func (c *Coll) Payload(payload ...interface{}) *Coll {
	route := c.routes[c.currentroute]
	jsonBytes, err := json.Marshal(payload[0])
	route.Body.JSON = string(jsonBytes)
	if err != nil {
		log.Printf("payload failed to marshal %s", err)
	}

	if payload[1] == nil || payload[1] == "" {
		route.Body.Mode = c.mode
	} else {
		route.Body.Mode = payload[1].(string)
	}

	c.routes[c.currentroute] = route
	return c
}
