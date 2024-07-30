# Messagio


## Run local
```sh
docker-compose up --build
```
```sh
localhost:7070/send_message
```
```sh
localhost:7070/get_stats
```

## Remote address
[http://193.168.173.123/send_message](http://193.168.173.123/send_message)
[http://193.168.173.123/get_stats](http://193.168.173.123/get_stats)

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
```json
{
    "processed_count": 11,
    "notProcessed_count": 0,
    "total_count": 11
}
```  
