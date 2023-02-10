# core

## System components

* __Executor__ - The server which traverses the workflow's DAG, evaluates traversal rules, and executes/calls workflow processes.
* __Reducer__ - The server which listens on Orchid's event bus and evaluates state updater rules.
* __Scheduler__ - The service which handles requests to resume workflow execution to the next available executor node. These requests may come from the reducer or from an external process previously called by the executor.
* __REST API__ - An administrative API used to manage and monitor the state of workflows.
