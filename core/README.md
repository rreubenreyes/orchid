# core

## System components

* __Executor__ - The server which traverses the workflow's DAG, evaluates traversal rules, and executes/calls workflow processes.
* __Reducer__ - The server which listens on Orchid's event bus and evaluates state updater rules.
* __Scheduler__ - The worker which handles requests to resume workflow execution to the next available executor node. These requests may come from the reducer or from an external process previously called by the executor.
* __gRPC API__ - Used by:
  * The reducer and external processes, to communicate requests to the scheduler.
  * All three other components, to publish updates to workflow state.
