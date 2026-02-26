# Program Browser

A tool to list programs on bug bounty platforms.

It fetches the program and formats them to a reusable format including the full
parsed scope to automate recon afterwards.

## Capabilities

For now, this program only supports bugcrowd public programs.

## Usage

The program must be used with a config file. It is used for output of the
program.

Here is an example config using redis

```yaml
# config.yaml
output:
    redis:
        address: localhost:6379 # the address of redis
        db: 0 # the redis db to use
        name: programs # the name of the queue to push to
        password: abc # The password for redis. If not set, `REDIS_PASSWORD` env variable is used
```

The program can be ran using the following command :

```shell
go run ./cmd/browser/ -config config.yaml
```

Here is a full config with commented parts :

```yaml
output:
  # redis:
  #   address: localhost:6379
  #   db: 0
  #   name: foundPrograms
  #   password: adsfa
  file:
    format: json
    # filename: test.json

input:
  filters:
  - glob: '*program1*' # filter by globbing
    insensitive: true  # case insensitive check

  - regex: '.*program[0-9].*' # filter by regex
    # insensitive: false

  - exact: 'Program1' # filter by exact name
    insensitive: true # case insensitive check

  # add extra entries for custom programs

  # extraEntries:
  # - platform: custom
  #   platform_id: 0110-0111-0110-0111
  #   scope:
  #     allowed_endpoints:
  #     - scheme: https
  #       host:
  #         suffix: example.com
  #         wildcard: true
  #     denied_endpoints: []
  #   name: Moovit Managed Bug Bounty Program
  #   url: ""

  # Config for bugcrowd
  bugcrowd:
    # Enable bugcrowd provider
    enable: true
```

### Outputs 

Everything in this section must be configured under the `output` section in
the config file

#### Redis

The found programs are outputed to a redis queue with the following config.

| name | type | default value | description |
| ---- | ---- | ----          | ---         |
| `redis.address` | `string` | `''` | The address of the redis database |
| `redis.db` | `int` | `0` | The db number of redis |
| `redis.name` | `string` | `''` | The name of the queue to output to |
| `redis.password` | `string` | `REDIS_PASSWORD` environment variable | The password to connect to redis database |

#### File

The found programs are outputed to a file or stdout (by default) using a given
format

| name | type | default value | description |
| ---- | ---- | ----          | ---         |
| `file.format` | `string` | `''` | The format to print to accepted values are `json` and `yaml` |
| `file.filename` | `string` | `stdout` | The file to output to. |
