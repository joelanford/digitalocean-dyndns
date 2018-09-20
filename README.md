# digitalocean-dyndns
Automatically keep your DigitalOcean DNS records up-to-date

## Docker
```
docker pull joelanford/digitalocean-dyndns
```

## Usage
```
$ ./digitalocean-dyndns help
NAME:
   digitalocean-dyndns - Automatically keep your DigitalOcean DNS records up-to-date

USAGE:
   digitalocean-dyndns [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token value, -t value     DigitalOcean API access token [$DIGITALOCEAN_TOKEN]
   --domain value, -d value    Domain to update [$DIGITALOCEAN_DOMAIN]
   --name value, -n value      Record name to update [$DIGITALOCEAN_NAME]
   --interval value, -i value  Update interval (default: 1h0m0s) [$DIGITALOCEAN_INTERVAL]
   --help, -h                  show help
   --version, -v               print the version
```
