# glitchtipctl

- A Commandline Tool for Glitchtip Error Tracking software written in Go.

- To start an instance of glitchtip ruuning locally, perform the following command, updating the docker-compose file with your email address and password.

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

## Usage

```bash
Usage:
  glitchtipctl [command]

Available Commands:
  completion         Generate the autocompletion script for the specified shell
  createOrganization Create a new organization using the GlitchTip API
  createProject      A brief description of your command
  createTeam         Create a new team within an organization using the GlitchTip API
  deleteOrganization A brief description of your command
  deleteProject      A brief description of your command
  deleteTeam         A brief description of your command
  deleteUser         A brief description of your command
  getMembers         A brief description of your command
  getOrganizations   List all organizations
  getProjects        Get a list of projects from your organization
  getTeams           Get a list of teams from your organization
  getUsers           A brief description of your command
  help               Help about any command
  organizations      A brief description of your command
  projects           A brief description of your command

Flags:
  -h, --help     help for glitchtipctl
  -t, --toggle   To toggle the debug mode

```
