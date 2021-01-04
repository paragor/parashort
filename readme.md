Мне нужен был сокращатель ссылок. Дешево, сердито, пнх.

```bash
make run
make run

HOST=http://localhost:8000
curl -X POST $HOST/api/v1/save --data '{"text":"huiiii"}' -H 'content-type: application/json'
curl $HOST/api/v1/list | jq '.list[]' -r | xargs -I{} curl $HOST/api/v1/url/{}
```
