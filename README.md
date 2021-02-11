# IoT Logger

This is a very basic, graphql based data logger server, designed to be used for IoT devices which cannot be debugged via serial, written in [Go](https://golang.org/).

## Installation

```shell
go install github.com/llamadeus/iot-logger@latest
```

## 

## Configuration
```shell
cp .env.example .env
```

Then edit `.env` and adjust it as desired.

**Alternatively, you can provide the variables defined in `.env.example` via standard environment variables.**

## Running the server
```shell
go run main.go
```

If `APP_ENV` is set to `development` you can test the server in your browser.
Just visit <http://localhost:8080/graphql> and try one of the following queries and mutations. 

## Schema

The server provides a graphql endpoint with the following schema:

```graphql endpoint
scalar Time

type Message {
    id: ID!
    text: String!
    timestamp: Time!
}

type Query {
    # Send a ping to the server.
    # She will reply with a "pong".
    ping: String!
}

type Mutation {
    # Add the given message to a channel.
    # If the mutation returns true you know that
    # at least one subscriber received the message.
    addMessage(channel: String!, message: String!): Boolean!
}

type Subscription {
    # Listens for messages in the given channel.
    messageAdded(channel: String!): Message!
}
```
