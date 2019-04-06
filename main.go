package main

import (
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	"github.com/rdegges/go-ipify"
)

func main() {
	// Discover the public IP.
	ip, err := ipify.GetIp()
	if err != nil {
		err = errors.Wrap(err, "find public IP error")
		log.Fatal(err)
	}

	// Initialize the Cloudflare API.
	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		err = errors.Wrap(err, "Cloudflare initialization error")
		log.Fatal(err)
	}

	// Update the IP.
	var (
		zoneID   = os.Getenv("CF_ZONE_ID")
		zoneType = os.Getenv("CF_ZONE_TYPE")
		zoneName = os.Getenv("CF_ZONE_NAME")
	)
	if err := updateIP(ip, zoneID, zoneType, zoneName, api); err != nil {
		err = errors.Wrap(err, "updateIP error")
		log.Fatal(err)
	}
}

func updateIP(ip, zoneID, zoneType, zoneName string, api *cloudflare.API) error {
	// Set the record to be found.
	record := cloudflare.DNSRecord{
		Type: zoneType,
		Name: zoneName,
	}

	// Fetch all records with a given name.
	records, err := api.DNSRecords(zoneID, record)
	if err != nil {
		return errors.Wrap(err, "fetch DNS records error")
	}
	if len(records) == 0 {
		return nil
	}

	// If the record already points to the current IP, we can safely return.
	if records[0].Content == ip {
		return nil
	}

	// Set the new IP at the record and send the change to Cloudflare.
	record.Content = ip
	if err := api.UpdateDNSRecord(zoneID, records[0].ID, record); err != nil {
		return errors.Wrap(err, "update DNS record error")
	}

	return nil
}
