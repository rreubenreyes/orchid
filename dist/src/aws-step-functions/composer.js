"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/**
 * Core rules for handling input/output processing:
 *
 * Any output from any state is persisted
 * The output of any state is either equivalent to or the superset of the input
 */
// function resolver(input: SFN.SerializableObject, opts: {
//     statePath: string;
// }): unknown {
//     const data = input;
//     const statePath = opts.statePath;
//     function getOutputOf(state,
// }
/**
 * combining pass state and task state below to illustrate what I'm trying to do here
 *
 * {
 *  "FirstState": {
 *      "Type": "Task"
 *      "Resource": "arn:aws:states:::lambda:FetchAppUser.invoke"
 *      "Result": {
 *          "appuser": {
 *              "humanId": "poop"
 *          }
 *      },
 *      "ResultSelector": {
 *          "human_id": "$.appuser.humanId"
 *      }
 *      "ResultPath": "$.data.FirstState.lambda.SomeFunction.result"
 *  }
 * }
 *
 * input to second state should be:
 * {
 *  "input": "whatever was passed in the input"
 *  "data": {
 *      "FirstState": {
 *          "lambda": {
 *              "FetchAppUser": {
 *                  "human_id": "poop"
 *              }
 *          }
 *      }
 *  }
 *
 * in order to access that value from the second state:
 *
 * "Parameters": {
 *  "human_id.$": "$.data.FirstState.lambda.SomeFunction.result.human_id"
 * }
 *
 * bottom line:
 *
 * - the developer should NEVER have to write a ResultPath
 * - the developer should NEVER and CAN NEVER write an OutputPath
 * - the developer will only have to worry about Parameters and ResultSelector
 *
 * some ideas:
 *
 * const userDataFetched = new State.Task('AppUserDataFetched', {
 *      resource: invokeLambda({ arn: arn:aws:us-west-2:...:FetchAppUser }),
 *      output: (lambdaResult) => ({
 *          'human_id': lambdaResult.get('appuser.humanId')
 *      }),
 * })
 *
 * const doSomethingWithUserData = new State.Task('DoSomething', {
 *      resource: invokeLambda({ arn: arn:aws:us-west-2:...:AnotherLambda }),
 *      parameters: (context) => ({
 *          'human_id': context.getState('AppUserDataFetched').getResult('human_id')
 *      }),
 * })
 */
