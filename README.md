# orchid

Orchid is an orchestration platform which allows developers to manage long-lived workflows using _federated_ processes.

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
