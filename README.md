TODO: update this

### Dev setup

- create api/.env.local
  - ```cp api/.env api/.env.local```
  - replace values with correct ones
- create dash/.env.local
    - ```cp dash/.env dash/.env.local```
    - replace values with correct ones
- create test deploy files
    - ```echo '{"commit": "testapi", "timestamp": "testapi"}' > api/tmp/deploy.json```
    - ```echo '{"commit": "testdash", "timestamp": "testdash"}' > dash/public/deploy.json```
- build and run the containers
  - ```docker compose up -d --build```
- add test user
  - ```docker exec -it murky_api bash -c 'go run cmd/addtestuser/main.go -db db/db.sqlite3'```

To access the site and api from mobile devices, add dev machine ip to .env files instead of "localhost":
- example api address: ```http://192.168.1.197:8080```
- example front address: ```http://192.168.1.197:5137```


### Staging setup

To test the deployments there is a local 'server' dir that the deploy scripts target.
It mirrors the prod server structure.

To get it ready for deploy scripts:
- create a .env.local there
  - ```cp dash/.env dash/.env.local```
  - replace values with correct ones
- set up a local Caddy server
  - install Caddy
  - example config is at server/caddy-config-example/Caddyfile
  - run it from project root
    - ```caddy run --config server/caddy-config-example/Caddyfile```

Then run the staging deploy scripts to build and deploy to the local server dir:
- ```bash deploy/build_staging.sh```
- ```bash deploy/deploy_staging.sh```


### Prod setup

Set up the prod server same way as the local server dir:
- create .env.local
- run Caddy

In addition, prod deploy requires ssh access to server with alias set to 'murky'.

Then run the prod deploy scripts to deploy to prod server:
- ```bash deploy/build_prod.sh```
- ```bash deploy/deploy_prod.sh```
