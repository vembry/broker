# broker

This is an attempt to create a modular message-broker using golang. The modularity: 
1. by running the broker as an domain which can be part of an existing service/monolith/microservice
2. by running the broker as a stand alone message broker which can exposed using gRPC
3. giving option back to user to specify storage/persistence they want to use for backup/restore during downtime/deployment

## structure

![alt text](image.png)

## flow

### enqueue
enqueue a message into a queue

```mermaid

flowchart TD

    0@{ shape: sm-circ, label: "Small start" } --> 1[enqueue]
    1 --> 2[find queue]
    2 --> 3{exists?}

    3 -- no --> 4[create new queue]
    4 --> 5

    3 -- yes --> 5[lock queue]

    5 --> 6[add new entry to the queue]
    6 --> 7[unlock queue]

    7 -->  8@{ shape: framed-circle, label: "Stop" }

```

### poll
poll a message from a queue

```mermaid

flowchart TD
    0a@{ shape: sm-circ, label: "Small start" } --> 1[poll]
    1 --> 2[find queue]
    2 --> 3{exists?}
    3 -- no --> 0b@{ shape: framed-circle, label: "Stop" }
    3 --> 4[lock queue]
    4 --> 5[extract out queue's task at '0' index]
    5 --> 6[put it on active message map]
    6 --> 7[unlock queue]
    7 --> 0b

```

### complete-poll
completes an active/polled message journey

```mermaid

flowchart TD
    0a@{ shape: sm-circ, label: "Small start" } --> 1[complete-poll]
    1 --> 2[delete active message]
    2 --> 0b@{ shape: framed-circle, label: "Stop" }


```

### sweeper
check for expiring message and put them back into the queue. `sweeper` will run in the background concurrently and currently set to run every 1 second

```mermaid

flowchart TD
    0a@{ shape: sm-circ, label: "Small start" } --> 1[iterate through the active message map]
    1 --> 2{has active message?}
    2 -- no --> 2.1[wait for interval]
    2.1 --> 1
    2 -- yes --> 3[check message expiry]
    3 --> 4{expired?}
    4 -- yes --> 5[extract it from active message map]
    5 --> 6[put it back to idleQueue map]
    6 --> 1

```