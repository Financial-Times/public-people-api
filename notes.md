
Shell commands
--------------

    go install && public-people-api --neo-url http://localhost:18080/__neo4j-red # connect to neo4j in XP via an SSH LocalForward tunnel.


curl commands
-------------

    curl -v http://localhost:8080/people/9b3fca66-a028-468e-8a77-d8259d621ff7 # should return 200
    curl -v http://localhost:8080/people/41038da6-90db-4953-950b-89cf2a0a432a # should return 301
