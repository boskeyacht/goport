# Goport

Goport is a golang implementation of a [Seaport gossip node](https://github.com/ProjectOpenSea/seaport-gossip).

## Overview

Currently, it connects & retrieves orders from other nodes in the network, subscribes to order events from the Seaport contract, and writes them both to a SQLite database.

>
> This project is still a work in progress. (Contributions are welcome!)
>


## Install

## Usage

- Create a `.env` file in `cmd/goport`.
- Navigate to `cmd/goport` and run `go run .`