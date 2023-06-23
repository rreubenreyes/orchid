# orchid

Orchid is an orchestration platform which allows developers to manage long-lived workflows using distributed processes.

## Design

TL;DR:

```
           ┌───────┐
     ┌────►│ State ├─────┐
     │     └───────┘     │
     │                   │
     │                   │
     │                   ▼
┌────┴─────┐       ┌───────────┐
│  Events  │◄──────┤ Processes │
└──────────┘       └───────────┘
```

* Workflows expressed as [DAGs](#dags)
    * Each node in the DAG can be associated with a process, which does _something_
        * "Internal" processes are executed by the workflow evaluator host, and logic is expressed directly in the DAG definition
        * "External" process are executed by some external process which interfaces with the DAG
* Workflows persist [immutable state](#state)
    * While state cannot be mutated, new iterations of the state are generated each time an event must effect the state
    * State can be represented as a doubly linked list
    * Each iteration of a workflow references one state snapshot
* Unidirectional data flow
    * Processes can dispatch events
    * Events are persisted in a log
    * As events are written to the log, a "state updater" receives those events, derives the correct state update from the event payload, and produces a state snapshot. 

## Appendix

### DAGs

In Orchid, __workflows__ are defined as _[directed acyclic graphs](https://en.wikipedia.org/wiki/Directed_acyclic_graph)_ (DAGs) whose nodes consist of any number of processes. For example, we can create a workflow that describes how to prepare a bowl of cereal:

```
                         ┌─────┐
     ┌─────────────┬─────┤Start├──┬─────────┐
     │             │     └─────┘  │         │
     │             │              │         │
     ▼             ▼              ▼         ▼
┌──────────┐  ┌───────────┐  ┌─────────┐  ┌───┐
│Place Bowl│  │Pour Cereal│  │Pour Milk│  │End│
└────┬─────┘  └────┬──────┘  └────┬────┘  └───┘
     │             │              │
     │             │              │
     │             ▼              │
     │           ┌────┐           │
     └──────────►│Wait│◄──────────┘
                 └────┘
```

### State

```json
{
  "state": {
    "bowl_placed": {
      "type": "boolean",
      "initial": false
    },
    "cereal_poured": {
      "type": "boolean",
      "initial": false
    },
    "milk_poured": {
      "type": "boolean",
      "initial": false
    }
  }
}
```