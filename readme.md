Мне нужен был сокращатель ссылок. Дешево, сердито, пнх.

```bash
make run
make run

SERVER_HOST=http://localhost:8000
curl -X POST $SERVER_HOST/api/v1/save --data '{"text":"huiiii"}' -H 'content-type: application/json'
curl $SERVER_HOST/api/v1/list | jq '.list[]' -r | xargs -I{} curl $SERVER_HOST/api/v1/url/{}
```
