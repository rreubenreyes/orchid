# orchid

Orchid is an orchestration platform which allows developers to manage long-lived workflows using federated processes.

## Design

The core of Orchid's design is as follows:

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

In Orchid, _workflows_ are defined as _[directed acyclic graphs](https://en.wikipedia.org/wiki/Directed_acyclic_graph)_ (DAGs) whose nodes consist of any number of _processes_. For example, we can create a workflow that describes how to prepare a bowl of cereal:

```
                         ┌─────┐
     ┌─────────────┬─────┤Start├──┬─────────┐
     │             │     └─────┘  │         │
     │             │              │         │
     ▼             ▼              ▼         ▼
┌──────────┐  ┌───────────┐  ┌─────────┐  ┌───┐
│Place bowl│  │Pour cereal│  │Pour milk│  │End│
└────┬─────┘  └────┬──────┘  └────┬────┘  └───┘
     │             │              │
     │             │              │
     │             ▼              │
     │           ┌────┐           │
     └──────────►│Wait│◄──────────┘
                 └────┘
```

Notice that even though we're simply talking about cereal, the graph above doesn't perform the steps in _any_ order, as we'd normally expect! How, then, are we supposed to know what to do?

A core feature of Orchid is that workflows represent a group of processes that could occur at any given time, or in other words, given any _state_ of the world. Workflows interpret that state and use it to decide which path, or paths, to travel.

Importantly, every workflow in Orchid has an _initial_ state. Continuing our example, let's consider the following initial state:

```json
{
  "bowl_placed": false,
  "cereal_poured": false,
  "milk_poured": false
}
```

In order to figure out how to traverse this workflow, we define our workflow with a set of rules, such as the following, which tells our workflow how to transition out of the "start" node:

```json
[
  {
    "id": "start",
    "process": null,
    "rules": [
      {
        "$and": [
          { "$eq": [".bowl_placed", false] }
          { "$eq": [".cereal_poured", false] }
          { "$eq": [".milk_poured", false] }
        ],
        "next": "place_bowl"
      },
      {
        "$and": [
          { "$eq": [".bowl_placed", true] }
          { "$eq": [".cereal_poured", false] }
          { "$eq": [".milk_poured", false] }
        ],
        "next": "pour_cereal"
      },
      {
        "$and": [
          { "$eq": [".bowl_placed", true] }
          { "$eq": [".cereal_poured", true] }
          { "$eq": [".milk_poured", false] }
        ],
        "next": "pour_milk"
      },
      {
        "$and": [
          { "$eq": [".bowl_placed", true] }
          { "$eq": [".cereal_poured", true] }
          { "$eq": [".milk_poured", true] }
        ],
        "next": "end"
      },
    ]
  }
]
```

Because of our initial state, we'd know that the first traversal of our workflow would be:
```
                         ┌─────┐
     ┌───────────────────┤Start│
     │                   └─────┘
     │
     ▼
┌──────────┐
│Place bowl│
└──────────┘
```

A workflow in Orchid always starts at the `start` state, and ends at one of two states: `wait` or `end`.

(TODO: if this note is still here, this intro example isn't done)
