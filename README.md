# StepmaniaDB Backend
This is the backend for the StepmaniaDB project. It is a REST API written in Go. It interfaces with a Postgres DB containing metadata for Stepmania packs, songs, and charts.
StepmaniaDB Backend returns JSON data that can be consumed by a frontend. I created a frontend that can be used, but other frontends can be written as needed. Any application that can consume HTTP endpoints and work with JSNO data can use this backend.

# Runing the application

## Setting up the Postgres DB
- Create a Postgres DB Server and DB
- Use this CREATE SQL script to setup the initial tables you will need:
    - `./scripts/setup_db.sql`
- I recommend creating a user with read-only access to the DB, and using that user for the application

## Downloading the application
- You can either:
    - Clone this repo, build it, optionally dockerize it, and run it
    - Use the Docker image I have created
        - Image repository at:
    - Use a binary I have built
        - Binary releases at:
- If you don't know what you are doing, and don't know what docker is, I recommend going with the prebuilt binary  

### Authenticating with Postgres DB  
- There are two ways currently to authenticate:
    - A. Setting the DB credentials in environment variables
    - B. Using AWS Secrets Manager to get credentials (If you store the secrets in AWS)
*Option A*
- Set these Environment Variables:
    - `DB_HOST`
    - `DB_NAME`
    - `DB_USER`
    - `DB_PASS`
    - `DB_PORT`
*Option B*
- Set these Environment Variables:
    - `AWS_DB_PASSWORD_SECRET_NAME`
    - `AWS_REGION`
    - `AWS_ACCESS_KEY_ID`
    - `AWS_SECRET_ACCESS_KEY`
    - `DB_NAME`
- The AWS Secret Value will at least need to contain the keys:
    - username
    - password
    - host
    - port

- *Note*: If the credentials are set in both ways, Option B will take precedence

