# VK Film Library project

## How to start project
- Start command: docker compose up --build
- End command: docker compose down

## cmd

- CMD-directory contains main.go file.

## Internal

- Internal directory contains DB config. PostgreSQL was chosen as the DBMS.

## pkg

- pkg directory contains handlers, logger, middleware and postgresql packages
- These packages implemenr the application logic. JWT-token is used for the authorization system. This system contains two levels: one for the admin role and other for the user role. 

### docs

- docs directory contains all needed documentation for the swagger.

## General info

- Database was created with these tools: OpenServer, DBeaver (with postgres connection)
- Swagger documentation contains all models and handlers
- GinRouter is used for the routing paths
- Authorization takes place via middleware

