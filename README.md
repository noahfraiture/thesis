# Intro

Hello, this repo contains the [paper](./paper/Fraiture_16702000_2025.pdf) of my master thesis, that you can build with `nix run .#pdf-process` if you want to built it with the cover pages, or `typst compile main.typ` without the cover page.

Four scenarios were developed during this master thesis, [Satellite Siege](./scenario-satellite), [Hotel Daemon](./scenario-hotel), [Patch War](./scenario-patch) and [Zheng Hijack](./scenario-zheng). The first 3 scenarios use VM with quickemu, and rely on the network of the host machine, and the last scenarios use docker container managed with docker compose. Every scenario has script to help the build and deployment process, written in either Nix or Nushell.
