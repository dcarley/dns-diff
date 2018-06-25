# dns-diff

Compare DNS records from two different servers.

## Tests

[dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html) is used as a
fixture server. You can start it with:

    docker-compose up -d

Then run the tests with:

    go test -v ./...

If necessary, the hosts and ports can be overridden with the following
environment variables:

- `DNS_PRI_HOST`
- `DNS_PRI_PORT`
- `DNS_SEC_HOST`
- `DNS_SEC_PORT`
