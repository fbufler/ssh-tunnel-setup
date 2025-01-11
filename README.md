# SSH Tunnel Setup

This tool allows you to create and manage a ssh tunnel setup.
You basically have two clients and a server. The goal is that one client can connect to the other client through the server by using the ssh tunnel.

## Installation

```bash
./scripts/build.sh <os> <arch>
```

## Usage

Configuration is supported in multiple ways:
- config.yaml file
- command line arguments

### Server

To prepare the server you need to run the following command:

```bash
ssh-tunnel-setup server
```

For configuration options see the [example config file server section](config.example.yaml) or run `ssh-tunnel-setup server --help`.

### Client

To prepare the client(s) you need to run the following commands:

For the connecting client:
```bash
ssh-tunnel-setup client
```

For the target client:
```bash
ssh-tunnel-setup target
```

For configuration options see the [example config file client section](config.example.yaml) or run `ssh-tunnel-setup client --help` or `ssh-tunnel-setup target --help`.

## Prerequisites

- A ssh server must be installed configured and running on the server.
- The clients must have access to the server.
- The clients must have the ssh client installed.
- The clients must have the ssh server installed if they are the target client.



