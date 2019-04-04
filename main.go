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

	// Fetch the DNS zone ID.
	zoneID := os.Getenv("CF_ZONE_ID")
	if zoneID == "" {
		log.Fatal("missing 'CF_ZONE_ID'")
	}

	// Fetch the DNS zone type.
	zoneType := os.Getenv("CF_ZONE_TYPE")
	if zoneType == "" {
		log.Fatal("missing 'CF_ZONE_TYPE'")
	}

	// Fetch the DNS zone name.
	zoneName := os.Getenv("CF_ZONE_NAME")
	if zoneName == "" {
		log.Fatal("missing 'CF_ZONE_NAME'")
	}

	// Update the IP.
	if err := updateIP(ip, zoneID, zoneType, zoneName, api); err != nil {
		err = errors.Wrap(err, "updateIP error")
		log.Fatal(err)
	}
}

func updateIP(ip, zoneID, zoneType, zoneName string, api *cloudflare.API) error {
	// Set the record to be find.
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

	record.Content = ip
	if err := api.UpdateDNSRecord(zoneID, records[0].ID, record); err != nil {
		return errors.Wrap(err, "update DNS record error")
	}

	return nil
}
