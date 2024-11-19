# Instabug Challenge

----------------------------


## INSTALLATION
Make sure you have docker installed on your machine.
```bash
    docker-compose up --scale instabug-chat=3
```
> **Note:** As Ruby on Rails is a fully featured framework, I used a mini framework in GoLang to build the task faster.

## USAGE
- Attached is a postman collection that you can use to test the API. in the top level directory of the project

## SYSTEM DESIGN

### Components:
- **NGINX** - Reverse proxy server that loads balance requests to the API.
- **Go API** - Handles requests
- **Elasticsearch** - Used to index the messages for search functionality.
- **MySQL** - Main datastore for the application.
- **Redis** - Used for caching the `chats` and `messages` data, it caches the GET list only [for simplicity]
- **RabbitMQ** - Used for handling the messages creation between the API and the Worker.
- **Go Worker** - Worker that handles the message creation.
- **Event/Listener** - Used to handle the message creation event to store indexed data in Elasticsearch.

### API Endpoints:
```bash
GET /applications
POST /applications
GET /applications/:Token
PUT /applications/:Token

POST /applications/:Token/chats
GET /applications/:Token/chats
GET /applications/:Token/chats/:ChatNumber

POST /applications/:Token/chats/:ChatNumber/messages
GET /applications/:Token/chats/:ChatNumber/messages
GET /applications/:Token/chats/:ChatNumber/messages/:MessageNumber
PUT /applications/:Token/chats/:ChatNumber/messages/:MessageNumber
```

### Decisions:
- Preventing collisions of `chat_number` per application and `message_number` per chat by using **Compound Indexes** in the database.

- Handling race conditions in getting the `chat_number` and `message_number` by using select for update in transactions.

### Project Structure:
> **Note:**  All system components are put together in the same structure for simplicity.
- **api** - Go API, you can find the routes, controllers, services, and models.
- **worker** - Go Worker, you can find the worker that handles the message creation in `./pkg/rabbitmq/message_worker.go`
- **postman-collection** - Postman collection for the API in `./ChatSys.postman_collection.json`
