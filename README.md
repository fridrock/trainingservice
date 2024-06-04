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
```