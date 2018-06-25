# dns-diff

Compare DNS records from two different nameservers. Useful for checking
consistency when migrating zones.

## Usage

Build and install with Go:

    go get -u github.com/dcarley/dns-diff

Run:

    $ dns-diff -pri 0.0.0.0:10053 -sec 0.0.0.0:20053 <<EOF
    same-a.example.com A
    value-a.example.com A
    EOF
    INFO[0000] ✔ same-a.example.com A
    WARN[0000] ✘ value-a.example.com A
    WARN[0000] - value-a.example.com.       60      IN      A       1.1.1.1
    WARN[0000] - value-a.example.com.       60      IN      A       2.2.2.2
    WARN[0000] + value-a.example.com.       60      IN      A       3.3.3.3
    WARN[0000] + value-a.example.com.       60      IN      A       4.4.4.4

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
