# polaris-service

Build: make build

docker build -t polaris-service .

# Running API server from command line
DBUSER=root DBPASSWD='Intel123!' DBHOST=localhost ./polaris

# Getting logs
docker logs polaris-polaris-mysql-1
docker logs polaris-polaris-service-1
