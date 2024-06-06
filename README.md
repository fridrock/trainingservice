# Training service
- [Exercise Groups](#exercise-groups)
## Exercise Groups
- EXCHANGE: sport_bot
#### CREATE
- ROUTING_KEY: trainings.exgroup.create
- REQUEST BODY:
```json
{
    "user_id": 2,
    "name": "Back"
}
```
- RESPONSE:
    - ROUTING_KEY: tgbot.exgroup.create
```text
SUCCESS: id: 2 was saved
ERROR: wrong input
ERROR: internal server error: error description
```
#### DELETE
- ROUTING_KEY: trainings.exgroup.delete
- REQUEST BODY:
```json
{
    "user_id": 2,
    "name": "Back"
}
```
- RESPONSE:
    - ROUTING_KEY: tgbot.exgroup.delete
```text
SUCCESS
ERROR: wrong input
ERROR: no rows deleted
```
#### FIND BY NAME
- ROUTING_KEY: trainings.exgroup.find
- REQUEST BODY:
```json
{
    "user_id": 2,
    "name": "Back"
}
```
- RESPONSE:
    - ROUTING_KEY: tgbot.exgroup.find
```text
SUCCESS: {"id":1,"user_id":2,"name":"Back"}
ERROR: wrong input
ERROR: sql.ErrNoRows
```