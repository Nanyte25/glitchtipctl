# glitchtipctl

- A Commandline Tool for Glitchtip Error Tracking software written in Go.

- To start an instance of glitchtip ruuning local perform the following command, updating the docker-compose file with your email address and password.

`docker-compose up -d`

- exec into the running webapp container to create and backend `admin` account if neccessary.


```
sudo docker ps //to get the container ID
sudo docker exec -it 7AAAAAAAAA bash

```
- Once in the container run the following to setup a backend `admin` user account in djano-admin.

```
./manage.py createsuperuser // its prompts for email address for account

```

