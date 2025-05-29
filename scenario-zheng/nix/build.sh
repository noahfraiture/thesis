#!/bin/sh
nix build .#zheng-victim
docker load -i result
nix build .#zheng-host
docker load -i result
nix build .#zheng-c2
docker load -i result
rm result
