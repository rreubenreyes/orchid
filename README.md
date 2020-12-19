# Orchid (working name)

This is a TypeScript/JavaScript library which aims to simplify working with AWS Step Functions. This is a work in progress.

## Goals

- Limit the verbosity of constructing a state machine
- Remove the burden of deciding how to name/namespace data that is passed around in a state machine
- Allow state machines to be constructed dynamically, and allow the logic to do so exist alongside the rest of users' code.
- Enable patterns for building reusable states/components

## Example

```javascript
const { StateMachine, States, Resources } = require('orchid');

/* Create a new Task state that gets an item from DynamoDB */
const userDataFetched = new States.Task('UserDataFetched', {
    task: Resources.dynamoDb.getItem({
        tableName: 'users'
        parameters: (context) => ({
            Key: { // This library does not provide a wrapper over the DynamoDB API (or any non-Step Functions APIs)!
                user_id: { S: context.execution.getInput('user_id') } // Get the data that was input to the workflow's execution
            }
        }),
    }),
    output: (result) => ({ // Output the result of the DynamoDB query into state machine data
        user_id: result.get('Item.user_id.S'),
        first_name: result.get('Item.first_name.S'),
        last_name: result.get('Item.last_name.S'),
        is_active: result.get('Item.is_active.BOOL'),
    })
});

const emailSentToUser = new States.Task('EmailSentToUser', {
    task: Resources.lambda.invoke({
        functionName: 'send_email',
        parameters: (context) => ({ // Access data that was output by any previous state
            first_name: context.getState(userDataFetched).getOutput('first_name'),
            last_name: context.getState(userDataFetched).getOutput('last_name'),
        })
    })
    output: (result) => ({
        success: result.get('success')
    })
});

const success = new States.Succeed('Success');

userDataFetched.setDownstream(emailSentToUser); // The next state after UserDataFetched will be EmailSentToUser
emailSentToUser.setDownstream(success);
success.terminal();

/* Create a new state machine from the states you've just defined */
const stateMachine = StateMachine.create({
    startAt: userDataFetched,
    states: [userDataFetched, emailSentToUser, success]
});

/* Get the state machine as JSON */
const stateMachineJson = stateMachine
```

The JSON output of the above example is:

```json
{
    "StartAt": "UserDataFetched",
    "States": {
        "UserDataFetched": {
            "Type": "Task",
            "Resource": "arn:aws:states:::dynamodb:getItem",
            "Parameters": {
                "TableName": "users",
                "Key": {
                    "user_id.$": "$.input.user_id"
                }
            },
            "ResultPath": "$.data.UserDataFetched.dynamodb.getItem.result",
            "ResultSelector": {
                "user_id.$": "$.Item.user_id.S",
                "first_name.$": "$.Item.first_name.S",
                "last_name.$": "$.Item.last_name.S",
                "is_active.$": "$.Item.is_active.BOOL"
            },
            "Next": "EmailSentToUser"
        },
        "EmailSentToUser": {
            "Type": "Task",
            "Resource": "arn:aws:states:::lambda:invoke",
            "Parameters": {
                "FunctionName": "send_email",
                "Payload": {
                    "first_name.$": "$.data.UserDataFetched.dynamodb.getItem.result.first_name",
                    "last_name.$": "$.data.UserDataFetched.dynamodb.getItem.result.last_name"
                }
            },
            "ResultPath": "$.data.UserDataFetched.dynamodb.getItem.result",
            "ResultSelector": {
                "success.$": "$.success"
            },
            "Next": "Success"
        },
        "Success": {
            "Type": "Succeed"
        }
    }
}
```
