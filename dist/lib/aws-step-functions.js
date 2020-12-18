"use strict";
/* this project is meant to output jsons which are valid amazon states language files.
 *
 * See: https://states-language.net/spec.html
 */
Object.defineProperty(exports, "__esModule", { value: true });
/**
 * All state types that are supported by AWS Step Functions.
 *
 * https://states-language.net/spec.html#statetypes
 */
var StatesType;
(function (StatesType) {
    StatesType["Task"] = "Task";
    StatesType["Parallel"] = "Parallel";
    StatesType["Map"] = "Map";
    StatesType["Pass"] = "Pass";
    StatesType["Wait"] = "Wait";
    StatesType["Choice"] = "Choice";
    StatesType["Succeed"] = "Succeed";
    StatesType["Fail"] = "Fail";
})(StatesType = exports.StatesType || (exports.StatesType = {}));
/**
 * All error types that are specific to AWS Step Functions. (StateNode machines may throw errors
 * which are NOT specific to AWS Step Functions.)
 *
 * https://docs.aws.amazon.com/step-functions/latest/dg/concepts-error-handling.html#error-handling-error-representation
 */
var StatesError;
(function (StatesError) {
    StatesError["ALL"] = "States.ALL";
    StatesError["DataLimitExceeded"] = "States.DataLimitExceeded";
    StatesError["Runtime"] = "States.Runtime";
    StatesError["Timeout"] = "States.Timeout";
    StatesError["TaskFailed"] = "States.TaskFailed";
    StatesError["Permissions"] = "States.Permissions";
})(StatesError = exports.StatesError || (exports.StatesError = {}));
