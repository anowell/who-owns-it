who-owns-it
===========

A simple service for keeping track of component ownership


Getting started
---------------

Create data/teams.json

    [
      {
        "name": "Team Awesome",
        "email": "team-awesome@example.com:",
        "members" [
          "Anthony Nowell"
        ]
      }
    ]


Create data/projects.json

    [
      {
        "name": "who-owns-it",
        "team": "Team Awesome",
        "aliases": ["whoowns", "who-owns", "whoownsit", "ownership"]
      }
    ]

Start who-owns-it server in a docker container. This bind mounts the data directory into the container which doesn't "just work" for remote docker hosts (e.g. Boot2Doc - sorry Mac users). It also bind mounts the app directory and runs it via 'go run' so that you can make changes, and simply restart your docker container for them to take effect.
    
    $ bin/start

Now you can curl it:

    $ curl -s $DOCKER_HOST:8080/whoowns
    {"Project":"who-owns-it","Team":"Team Awesome","Email":"
    
    
Alternatively, you can run this outside of Docker
  
    $ DATA_DIR=`pwd`/data go run app/server.go
    $ curl localhost:8080/whoowns 
