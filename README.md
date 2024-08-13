AspacePublisher
=======

ApacePublisher is small server built with the Golang Echo framework that automates pushing metadata from ArchivesSpace to other services, currently ArchivesWest and OCLC. Future workflows will include OregonDigital and Alma.

### Usage

The ids listed as a parameter refer to the final component in the ArchivesSpace URI, e.g. in the case of a resource: /repositories/2/resources/3476, the id is 3476. 
Endpoints:
  - EAD endpoints pull EAD xml from ArchivesSpace and then posts to ArchivesWest
    - `/ead/convert/:id`
      - param: ArchivesSpace resource id as above
      - returns XML adapted for uploading manually to ArchivesWest
    - `/ead/validate/:id`
      - param: ArchivesSpace resource id as above
      - returns message of success or errors
    - `/ead/upload/:id`
      - param: ArchivesSpace resource id as above
      - returns message of success or errors
  - OCLC endpoints pull MARC XML from ArchivesSpace, posts to OCLC, then in the case of a new record, updates the resource.
    - `/oclc/:id`
      - param: ArchivesSpace resource id as above
      - returns message or errors formatted as json
    - `/oclc/validate/:id`
    -  param: ArchivesSpace resource id as above
    -  returns message formatted as json


### Local development

- clone repo
- set up the .env. See the docker-compose for current variables needed by the system to run all of the supported processes.

Running directly on local system
- install go
- run `go get <package>` for the packages in go.mod OR wait for the system to tell you what needs to be installed in the next step
- `go run main.go`

Docker
- in one terminal: `docker-compose up`
- in another terminal: `docker-compose exec server bash`
- then: `go run main.go`

NOTE for local development: connecting to the UO ArchivesSpace API requires being on the VPN.

### Staging:
Docker (required) 
as above, with one additional step, run: `register-app .` (generates docker-compose.staging.yml)
`docker compose -f docker-compose.staging.yml up server`

### Production:
1. `go build`
2. scp the executable to the server
3. if needed, edit /etc/aspace-pub.env
4. `sudo systemctl stop aspace-pub`
5. move executable to /usr/local/aspace-pub
6. `sudo systemctl start aspace-pub`

logging: `journalctl -fu aspace_publisher`

NOTE this service is only reachable from campus or while on the VPN.


