Мне нужен был сокращатель ссылок. Дешево, сердито, пнх.

```bash
make run
make run

curl -X POST localhost:8000/api/v1/save --data '{"text":"huiiii"}' -H 'content-type: application/json'
curl localhost:8000/api/v1/list | jq '.list[]' -r | xargs -I{} curl localhost:8000/api/v1/url/{}
```
