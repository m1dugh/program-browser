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

### Outputs 

#### Redis

The found programs are outputed to a redis queue with the following config.

| name | type | default value | description |
| ---- | ---- | ----          | ---         |
| `redis.address` | `string` | `''` | The address of the redis database |
| `redis.db` | `int` | `0` | The db number of redis |
| `redis.name` | `string` | `''` | The name of the queue to output to |
| `redis.password` | `string` | `REDIS_PASSWORD` environment variable | The password to connect to redis database |
