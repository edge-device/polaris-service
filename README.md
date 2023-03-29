# polaris-service

Build: make build or make docker



# Running API server from command line
DBUSER=root DBPASSWD='Intel123!' DBHOST=localhost ./polaris

# Getting logs
docker logs polaris-polaris-mysql-1
docker logs polaris-polaris-service-1


http://192.168.1.12:8000/v1/device/my_org/waiting_room


mysql -h localhost --protocol tcp -u polaris -p
DELETE FROM users WHERE user_id="kedge.management@gmail.com";

INSERT INTO users (user_id, created_at, last_login, firstname, lastname) VALUES('kedge.management@gmail.com', 1513615539, 0, 'Mike', 'Miller');

## Secrets

blah blah blah