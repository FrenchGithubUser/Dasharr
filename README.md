# Dasharr

 Dashboard of your indexers' usage

 > Note: this project is in very early stages, expect broken stuff and breaking changes, this is not ready for public use yet !

 > Don't open any issues regarding bugs, but I'm opened to feature suggestions/requests

 ![header](images/0.png)
 ![header](images/1.png)

 ## Quickstart

 ### Docker

- Copy [docker-compose.yml](./docker-compose.yml), edit the fields you want (most importantly the `API_KEY` env var).
- Run the container `docker-compose up -d`
- Go to the config volume and enable the indexers you want in `config.json`
- Visit the webui and configure the indexers that require configuration
- Restart the container