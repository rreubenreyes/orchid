# orchid

Orchid is an orchestration platform which allows developers to manage long-lived workflows using federated processes.

## Design

Unidirectional data flow is at the core of Orchid's design. If you're familiar with [Redux](https://redux.js.org/), then this term might be familiar to you. All Orchid workflows adhere to the following data flow:

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

Notice that even though we're simply talking about cereal, the graph above doesn't perform the steps in _any_ order, as we'd normally expect! How, then, are we supposed to know what to do?

A core feature of Orchid is that workflows represent a group of processes that could occur at any given time, or in other words, given any __state__ of the world. Workflows interpret that state and use it to decide which path, or paths, to travel.

Importantly, every workflow in Orchid has an _initial_ state. Continuing our example, let's consider the following initial state:

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

In order to figure out how to traverse this workflow, we define our workflow with a set of rules, such as the following, which tells our workflow how to transition out of the "start" node:

```json
{
  "nodes": [
    {
      "id": "start",
      "process": null,
      "rules": [
        {
          "rule": {
            "$and": [
              { "$eq": [".bowl_placed", false] },
              { "$eq": [".cereal_poured", false] },
              { "$eq": [".milk_poured", false] }
            ]
          },
          "next": "place_bowl"
        },
        {
          "rule": {
            "$and": [
              { "$eq": [".bowl_placed", true] },
              { "$eq": [".cereal_poured", false] },
              { "$eq": [".milk_poured", false] }
            ]
          },
          "next": "pour_cereal"
        },
        {
          "rule": {
            "$and": [
              { "$eq": [".bowl_placed", true] },
              { "$eq": [".cereal_poured", true] },
              { "$eq": [".milk_poured", false] }
            ]
          },
          "next": "pour_milk"
        },
        {
          "rule": {
            "$and": [
              { "$eq": [".bowl_placed", true] },
              { "$eq": [".cereal_poured", true] },
              { "$eq": [".milk_poured", true] }
            ]
          },
          "next": "end"
        }
      ]
    }
  ]
}
```

Because of our initial state, we know that the first traversal of our workflow would be:

```
                         ┌─────┐
     ┌───────────────────┤Start│
     │                   └─────┘
     │
     ▼
┌──────────┐
│Place Bowl│
└──────────┘
```

Now we're at the `place_bowl` node. Remember that we're making a bowl of cereal here; at this point, we want to actually _do_ something to make sure we make progress towards that goal. So, `place_bowl` will act as a _process node_. A __process__ in Orchid is the specific action that we take upon traveling to a particular node in our workflow. We can define processes in many ways, in any language we'd like, but the simplest way is via an _inline_ process directly in the workflow definition:

```json
{
  "nodes": [
    {
      "id": "place_bowl",
      "process": {
        "type": "inline",
        "emits": ["bowl_placed"],
        "definition": {
          "type": "emit_event",
          "event": {
            "type": "bowl_placed",
            "payload": null
          }
        }
      },
      "rules": [
        {
          "rule": "always"
          "next": "wait"
        }
      ]
    }
  ]
}
```

The simple process we define above emits an __event__ of type `bowl_placed`. Events are how processes in Orchid communicate their side effects to the rest of the workflow.

Importantly, __emitting the `bowl_placed` event _does not_ directly update the workflow's state__. In order to update the workflow state and effect a change in the next workflow iteration, we must write a state updater.

A __state updater__ is a special type of process which may never emit any events. Its sole purpose, as its name suggests, is to update workflow state in response to any events emitted by workflow processes. Like normal processes, state updaters can (and, for the most part, are) be written in any language using complex logic. However, it's still possible to write inline state updaters directly in the workflow definition. For demonstration purposes, we will do just that.

```json
{
  "updater": {
    "type": "inline",
    "definition": [
      {
        "event_type": "bowl_placed",
        "rules": [
          {
            "rule": {
              "$eq": ["state.bowl_placed", false],
            },
            "update": {
              "$set": ["state.bowl_placed", true]
            }
          }
        ]
      }
    ]
  }
}
```

(TODO: if this note is still here, this intro example isn't done)
