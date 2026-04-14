### Dev setup

- create .env file
  - ```cp .env.example .env```
  - replace values with correct ones
- create test deploy file
    - ```echo '{"commit": "testapi", "timestamp": "testapi"}' > tmp/deploy.json```
- build and run the containers
  - ```docker compose up -d --build```
- add test user
  - ```docker exec -it murky_api bash -c 'go run cmd/addtestuser/main.go -db db/db.sqlite3'```

To access the site and api from mobile devices, add dev machine ip to .env files instead of "localhost": ```http://192.168.1.197:8080```


### Prod setup

Deployment is done with GitLab ci. It runs the deployment script on main branch changes. Make sure to have necessary variables set in GitLab.

### TODO

- add music endpoints
  - ripping
  - listing
  - streaming
