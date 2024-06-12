# Training service
- [Exercise Groups](#exercise-groups)
- [Trainings](#trainings)
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
#### UPDATE BY NAME
- ROUTING_KEY: trainings.exgroup.update
- REQUEST BODY:
```json
{
    "user_id": 2,
    "name": "Back",
    "newname":"NewBack"
}
```
- RESPONSE:
    - ROUTING_KEY: tgbot.exgroup.update
```text
SUCCESS
ERROR: wrong input
ERROR: no rows updated 
```
#### FIND BY USER
- ROUTING_KEY: trainings.exgroup.findByUser
- REQUEST BODY:
```json
{
    "user_id":3
}
```
- RESPONSE:
    - ROUTING_KEY: tgbot.exgroup.findByUser
```text
SUCCESS: [
        {
        "id": 1,
        "user_id": 2,
        "name": "Back"
        },
        {
        "id": 2,
        "user_id": 2,
        "name": "Front"
        },
        {
        "id": 3,
        "user_id": 2,
        "name": "Chest"
        }
        ]
```
## Trainings
- EXCHANGE: sport_bot
#### START TRAINING
- ROUTING_KEY: trainings.training.start
- REQUEST BODY:
```json
{
    "user_id":1
}
```
- RESPONSE:
    - ROUTING_KEY: tgbot.training.start
```text
ERROR: wrong input
SUCCESS: id:12
```
#### FINISH TRAINING
- ROUTING_KEY: trainings.training.finish
- REQUEST BODY:
```json
{
    "user_id":1
}
```
- RESPONSE:
    - ROUTING_KEY: tgbot.training.finish
```text
ERROR: wrong input
ERROR: error finishing training: Empty non-finished trainings list
SUCCESS
```
#### GET TRAININGS
- ROUTING_KEY: trainings.training.get
- REQUEST BODY:
```json
{
    "user_id":1
}
```
- RESPONSE:
    - ROUTING_KEY: tgbot.training.get
```text
ERROR: wrong input
ERROR: error finishing training: Empty non-finished trainings list
SUCCESS: [
{
"id": 1,
"user_id": 2,
"begins": "2024-06-12T21:23:03.7097226+03:00",
"finish": "2024-06-12T21:23:03.7097226+03:00"
},
{
"id": 2,
"user_id": 2,
"begins": "2024-06-12T21:23:03.7097226+03:00",
"finish": "2024-06-12T21:23:03.7097226+03:00"
},
{
"id": 3,
"user_id": 2,
"begins": "2024-06-12T21:23:03.7097226+03:00",
"finish": "2024-06-12T21:23:03.7097226+03:00"
}
]
```