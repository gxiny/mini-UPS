# ERSS Final Project

This is the submission repository of ERSS final project.

Our task is to implement a simple version of
Useless Parcel Service (UPS).

## Configuration

You may want to change the secret key used by Django:
edit `secret.txt` and **do not** commit it to the Git repository.

## Build

The project requires using Docker.  To build the images, run

    # docker-compose build

If you change the source code, build the images again.

## Run

    # docker-compose up

Please note that there's a race condition in this setup,
especially when the database container is first started up.
Currently the backend program will wait until the database is up,
but Django may fail.  In this case, stop the containers and restart
them.

Or, start the database container first, wait for a while, and then
start the rest containers.

    # docker-compose up -d db # and wait for a while
    # docker-compose up
