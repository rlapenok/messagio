# Messagio


## Run
```sh
docker-compose up --build
```


## API
#### /send_message
- `POST` : Create a new message
```json
{
    "msg":"hello"
}
```

#### /get_stats
- `GET` : Get processed messages
