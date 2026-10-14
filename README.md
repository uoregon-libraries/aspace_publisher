AspacePublisher
=======

ApacePublisher is small server built with the Golang Echo framework that automates pushing metadata from ArchivesSpace to other services, currently ArchivesWest, OCLC, and Alma. Future workflows may include OregonDigital.

### Usage
Before starting any workflow, log in using Aspace credentials at /as/login.html. If your session has expired, it will notify you to log in again.
Most of the workflows are available at /launch.html. The form asks for the Aspace resource id and then you can choose from the following workflows:
- Validate EAD at Archives West
- Create/Update EAD at Archives West
- Validate MARC at OCLC
- Create/Update MARC at OCLC
- Create/Update Alma record at Ex Libris

After submitting the launch form, the system will generate a filename and display a link; this is a report which will hold any data the system records while running the workflow. The presence of the link does not necessarily indicate that the process has finished; most of the workflows should not take longer than a minute or two. This is NOT true of the Alma workflows; while most of the work should complete fairly quickly, if a new bib is being linked to the NZ, this will take awhile. To see whether or not an Alma workflow is still in progress, you can go to /riverui/jobs.

A separate workflow for creating digital objects can be started by uploading a json file at /as/do.html. See the wiki page for additional information regarding this workflow.

### Local development

- clone repo
- set up the .env. See the docker-compose for current variables needed by the system to run all of the supported processes.

Podman
- `podman build -t aspace_publisher:tag .`
- `podman build -t river_service:tag -f Docker_worker .`
- `podman compose up`

NOTE for local development: connecting to the UO ArchivesSpace API requires being on the VPN.

### Staging:
Staging is now using podman

### Production:
Production is currently using a version that is not containerized.
1. `go build`
2. scp the executable and any static files to the server
3. if needed, edit /etc/aspace-pub.env
4. `sudo systemctl stop aspace-pub`
5. move executable to /usr/local/aspace-pub
6. move static files to the appropriate place in /usr/local/aspace-pub
7. `sudo systemctl start aspace-pub`

logging: `journalctl -fu aspace-pub`

NOTE this service is only reachable from campus or while on the VPN.


