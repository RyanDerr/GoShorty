# GoShorty
URL Shortener Application

## Example Commands
These commands utilize the `shorten.sh` and `resolve.sh` scripts to interact with the URL shortener service.

### Shorten URL
To shorten a URL, you can use the following command:
```sh
./scripts/shorten.sh --url "http://example.com" --short "exmpl" --exp "1h
```
This will shorten the URL `http://example.com` with the custom short `exmpl` and an expiration time of 1 hour

### Resolve URL
To resolve a shortened URL, you can use the following command:
```sh
./scripts/resolve.sh --short exmpl
```
This will resolve the shortened URL `exmpl` and return the original URL.
