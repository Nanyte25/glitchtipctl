# glitchtipctl

A Commandline Tool for GlitchTip Error Tracking software written in Go.

## Prerequisites

1. You need to have a GlitchTip instance running locally or remotely.
2. You need an API token from GlitchTip to authenticate API requests.

## Getting Started

### Step 1: Start GlitchTip Locally

To start an instance of GlitchTip running locally, perform the following command after updating the `docker-compose.yaml` file with your email address and password:

```bash
docker-compose up -d
```

- Step 2: Create an Admin Account

If necessary, you can exec into the running webapp container to create a backend admin account.

    Find the container ID by running:

```bash

sudo docker ps  # to get the container ID
```

- Exec into the container using:

```bash

sudo docker exec -it <CONTAINER_ID> bash
```

- Once inside the container, create the admin user account by running the following command:

```bash

    ./manage.py createsuperuser  # This will prompt you for the email address and password for the account
```

- Step 3: Set Up the .env File

- Create a .env file in the root of your project directory. This file will store your GlitchTip API token securely.

- Here’s an example .env file:

```bash

GLITCHTIP_API_TOKEN=your-glitchtip-api-token-here
```

- Replace your-glitchtip-api-token-here with the actual API token from GlitchTip.
- Step 4: Install Dependencies

- Install the necessary Go dependencies, including the godotenv package to load environment variables:

```bash

go get github.com/joho/godotenv

```
- Step 5: Build the CLI Tool

- Once everything is set up, you can build the glitchtipctl binary by running:

```bash

go build -o glitchtipctl
```
- This will create a binary named glitchtipctl in the current directory.
- Step 6: Run the CLI Tool

- After building the tool, you can start using glitchtipctl to interact with your GlitchTip instance.
## Usage

-- The glitchtipctl tool provides various commands for managing organizations, projects, teams, users, and more.

```bash

Usage:
  glitchtipctl [command]

Available Commands:
  completion         Generate the autocompletion script for the specified shell
  createOrganization Create a new organization using the GlitchTip API
  createProject      Create a new project in GlitchTip
  createTeam         Create a new team within an organization using the GlitchTip API
  deleteOrganization Delete an organization
  deleteProject      Delete a project
  deleteTeam         Delete a team
  deleteUser         Delete a user
  getMembers         Fetch the members of an organization by organizational slug
  getOrganizations   List all organizations
  getProjects        Get a list of projects from your organization
  getTeams           Get a list of teams from your organization
  getUsers           Fetch the users of an organization
  help               Help about any command
  projects           Manage projects within an organization

Flags:
  -h, --help     help for glitchtipctl
  -t, --toggle   To toggle the debug mode

For more details about a command, use glitchtipctl [command] --help.
Example Commands
```
- Here are some example commands you can use to interact with your GlitchTip instance:
Create a New Organization

```bash

./glitchtipctl createOrganization --name "NewOrg" --slug "new-org"
````
- Create a New Project

```bash

./glitchtipctl createProject --name "My New App" --slug "my-new-app" --organization "org-slug" --team "team-slug" --platform "react"
```

## List All Organizations

```bash

./glitchtipctl getOrganizations
```
## List All Projects

```bash

./glitchtipctl getProjects
```
## List Members of an Organization

```bash

./glitchtipctl getMembers [organization_slug]
```
- Delete a Project

```bash

./glitchtipctl deleteProject --organization "org-slug" --slug "project-slug"

```
- Step 7: Contributing

- Feel free to open issues or submit pull requests if you want to contribute to the development of glitchtipctl. Contributions are welcome!