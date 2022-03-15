# Go gin demo
### HOW TO RUN ?
- Using docker-compose

```
docker-compose up -d --build
```
- Sample request
```
curl --request POST 'localhost:8000/user/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "test@email.com",
    "username": "test",
    "password": "12345678"
}'
```

- APIs
```azure
POST   /user/register 
POST   /user/login
GET    /user/:userId
GET    /user 
PATCH  /user/:userId/active-status
```