## Installation

`nats-subscribe` requires Go 1.11 or later.

```bash
$ go install github.com/adzeitor/nats-subscribe
```

## Usage

```bash
$ NATS_URL="nats://nats-server:4222" nats-subscribe subject1 subject2
```

###### NATS streaming

```bash
$ NATS_URL="nats://nats-server:4222" nats-subscribe -streaming subject1 subject2
```