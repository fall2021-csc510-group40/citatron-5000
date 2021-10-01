# Troubleshooting

This document describes how to troubleshoot issues in the project.

## Building

The project should build out of the box using Docker (or Docker compose).

However, if things need to be built individually:

- Core: change directories to `core` and then run: `docker build -t core .`
- Telegram bot: change directories to `telegram-bot` and then run: `docker build -t telegram-bot .`

## Deploying

Deployment should work out of the box via Docker compose (see bot installation in the README).

However, to troubleshoot any issues:

- Core: `docker logs core`
- Telegram bot: `docker logs telegram-bot`
- Mongo: `docker logs mongo`
- Mongo express: `docker logs mongo-express`