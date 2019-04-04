# DIP - Dynamic IP
## How to use it
Export these environment variables:
```shell
export CF_API_KEY=
export CF_API_EMAIL=
export CF_ZONE_ID=# this can be get at the domain page.
export CF_ZONE_TYPE=# the DNS record type.
export CF_ZONE_NAME=# the DNS record name.
```

Run the application:
```shell
go run main.go
```

Check the Cloudflare Admin page to see if the value has changed.