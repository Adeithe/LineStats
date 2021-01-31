# LineStats

A bot for Twitch that stores logs to provide user quotes and statistics.

## Getting Started

LineStats uses PostgreSQL for data storage and runs as a [Docker](https://docker.com/) container and needs to be built as an image for use.

Once you have installed [Docker](https://docker.com/) and have it running, building the image and getting the bot running is pretty simple.

`$ docker build -t linestats .`

Once the image has been built you can start the container using the following command. (Be sure to set up your environment variable)

`$ docker run -d --name=linestats --restart=always --env-file=./.env linestats`

If you're advanced and want to set up Prometheus metrics, just expose port 9091 on the container.
