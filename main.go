package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/digitalocean/godo"
	"github.com/revh/ipinfo"
	"github.com/urfave/cli"
)

var (
	ErrNoChange = fmt.Errorf("no change to external IP address")
)

func main() {
	app := cli.NewApp()

	app.Name = "digitalocean-dyndns"
	app.Usage = "Automatically keep your DigitalOcean DNS records up-to-date"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token, t",
			Usage:  "DigitalOcean API access token",
			EnvVar: "DIGITALOCEAN_TOKEN",
		},
		cli.StringFlag{
			Name:   "domain, d",
			Usage:  "Domain to update",
			EnvVar: "DIGITALOCEAN_DOMAIN",
		},
		cli.StringFlag{
			Name:   "name, n",
			Usage:  "Record name to update",
			EnvVar: "DIGITALOCEAN_NAME",
		},
		cli.DurationFlag{
			Name:   "interval, i",
			Usage:  "Update interval",
			Value:  time.Hour,
			EnvVar: "DIGITALOCEAN_INTERVAL",
		},
	}
	app.Action = func(c *cli.Context) error {
		token := c.GlobalString("token")
		domain := c.GlobalString("domain")
		name := c.GlobalString("name")
		interval := c.GlobalDuration("interval")

		if token == "" {
			return fmt.Errorf("token is required")
		}

		if domain == "" {
			return fmt.Errorf("domain is required")
		}

		if name == "" {
			return fmt.Errorf("name is required")
		}

		if interval < time.Second {
			return fmt.Errorf("interval must be at least 1s")
		}
		log.Printf("config: domain=%s, name=%s, interval=%s\n", domain, name, interval)

		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		oauthClient := oauth2.NewClient(context.Background(), tokenSource)
		client := godo.NewClient(oauthClient)

		if err := updateRecord(client, domain, name); err != nil {
			log.Printf("could not update record: %s", err)
		}

		ticker := time.NewTicker(interval)
		for range ticker.C {
			if err := updateRecord(client, domain, name); err != nil {
				log.Printf("could not update record: %s", err)
			}
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func updateRecord(client *godo.Client, domain, name string) error {
	// Get the current external IP
	ip, err := ipinfo.MyIP()
	if err != nil {
		return fmt.Errorf("could not get external IP")
	}

	// Get the list of existing records and iterate to find the record we want to update.
	records, resp, err := client.Domains.Records(context.Background(), domain, nil)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("could not retrieve existing records (domain=%s): %s", domain, err)
	}

	var record godo.DomainRecord
	var found bool
	for _, r := range records {
		if r.Name == name && r.Type == "A" {
			record = r
			found = true
		}
	}
	if !found {
		return fmt.Errorf("could not locate record (domain=%s, name=%s)", domain, name)
	}

	// If the record already has the correct IP, no need to send the edit request
	if record.Data == ip.IP {
		log.Printf("skipping update: no change detected (domain=%s, id=%d, name=%s, data=%s)", domain, record.ID, name, record.Data)
		return nil
	}

	// Update the domain record
	drer := &godo.DomainRecordEditRequest{
		Data: ip.IP,
	}
	_, resp, err = client.Domains.EditRecord(context.Background(), domain, record.ID, drer)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("could not update record (domain=%s, id=%d, name=%s, from=%s, to=%s): %s", domain, record.ID, name, record.Data, drer.Data, err)
	}
	log.Printf("updated domain record (domain=%s, id=%d, name=%s, from=%s, to=%s)", domain, record.ID, name, record.Data, drer.Data)

	return nil
}
