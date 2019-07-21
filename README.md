# graphql-populate

Docker image that populates database using scripts, all loaded via environment vars.

This program can be used to run GQL scripts with GraphQL mutations and queries.
It can be used to populate any GraphQL server. To do that you need to configurate
some env vars:

- **GRAPHQL_URL**: the url for the server to populate.

- **GQL_CHECK_TABLES**: a GQL query file, to check if desired tables exists. In
case they don't, the program will wait and check again.
- **GQL_POPULATE_FILE**: a GQL mutation file, with the mutation to populate the
database with GraphQL interface. An upsert mutation is more suitable, since it
can be rerun as many times as your docker image is restarted.
