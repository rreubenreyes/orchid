/* this project is meant to output jsons which are valid amazon states language files.
 *
 * See: https://states-language.net/spec.html
 */

/**
 * The union of all primitive JavaScript types which can be serialized as valid AWS Step Functions values.
 */
export type SerializablePrimitive = string | number | null;

/**
 * Recursive type. Fields can be SerializablePrimitives or the defined type itself.
 */
export interface SerializableObject {
    [index: string]: SerializablePrimitive | SerializableObject;
    [index: number]: SerializablePrimitive | SerializableObject;
}

export interface SerializableRecord {
    [index: string]: SerializablePrimitive | SerializableObject;
}

/**
 * The union of all JavaScript types which can be serialized as valid AWS Step Functions values.
 */
export type Serializable = SerializablePrimitive | SerializableObject;

/**
 * All state types that are supported by AWS Step Functions.
 *
 * https://states-language.net/spec.html#statetypes
 */
export enum StatesType {
    Task = 'Task',
    Parallel = 'Parallel',
    Map = 'Map',
    Pass = 'Pass',
    Wait = 'Wait',
    Choice = 'Choice',
    Succeed = 'Succeed',
    Fail = 'Fail'
}

/**
 * All error types that are specific to AWS Step Functions. (StateNode machines may throw errors
 * which are NOT specific to AWS Step Functions.)
 *
 * https://docs.aws.amazon.com/step-functions/latest/dg/concepts-error-handling.html#error-handling-error-representation
 */
export enum StatesError {
    ALL = 'States.ALL',
    DataLimitExceeded = 'States.DataLimitExceeded',
    Runtime = 'States.Runtime',
    Timeout = 'States.Timeout',
    TaskFailed = 'States.TaskFailed',
    Permissions = 'States.Permissions'
}

/**
export interface for a Step Functions state.
 *
 * https://states-language.net/spec.html#statetypes
 */
export interface StateNode {
    Type: StatesType;
    Comment?: string;
}

/**
 * Step Functions retrier. Used to implement retry logic.
 *
 * https://docs.aws.amazon.com/step-functions/latest/dg/concepts-error-handling.html#error-handling-retrying-after-an-error
 */
export interface Retrier {
    ErrorEquals: Array<StatesError | string>;
    IntervalSeconds?: number;
    MaxAttempts?: number;
    BackoffRate?: number;
}

/**
 * Step Functions catcher. Used to implement fallback logic.
 *
 * https://docs.aws.amazon.com/step-functions/latest/dg/concepts-error-handling.html#error-handling-fallback-states
 */
export interface Catcher {
    ErrorEquals: Array<StatesError | string>;
    Next: StateNode;
    ResultPath?: string;
}

/**
 * Step Functions Pass state. Commonly used to mock the presence of a state node.
 *
 * https://states-language.net/spec.html#pass-state
 */
export interface PassStateNode extends StateNode {
    Type: StatesType.Pass;
    Next?: string;
    End?: boolean;
    InputPath?: string;
    OutputPath?: string;
    ResultPath?: string;
    Result?: Serializable;
    Parameters?: Serializable;
}

/**
 * Step Functions Task state. Used to perform a scalar task.
 *
 * https://states-language.net/spec.html#task-state
 */
export interface TaskStateNode extends StateNode {
    Type: StatesType.Task;
    Next?: string;
    End?: boolean;
    InputPath?: string;
    OutputPath?: string;
    ResultPath?: string;
    Parameters?: Serializable;
    ResultSelector?: string | SerializableRecord;
    Retry?: Array<Retrier>;
    Catch?: Array<Catcher>;
    Resource: string;
    TimeoutSeconds?: number;
    TimeoutSecondsPath?: string;
    HeartbeatSeconds?: number;
    HeartbeatSecondsPath?: string;
}

/**
 * Rules which are valid in a Step Functions Choice state.
 *
 * https://states-language.net/spec.html#choice-state
 */
export interface ChoiceRule {
    And?: Array<ChoiceRule>;
    Or?: Array<ChoiceRule>;
    Not?: ChoiceRule;
    IsBoolean?: boolean;
    IsNull?: boolean;
    IsNumeric?: boolean;
    IsPresent?: boolean;
    IsString?: boolean;
    IsTimestamp?: boolean;
    NumericEquals?: number;
    NumericEqualsPath?: string;
    NumericGreaterThan?: number;
    NumericGreaterThanPath?: string;
    NumericGreaterThanEquals?: number;
    NumericGreaterThanEqualsPath?: string;
    NumericLessThan?: number;
    NumericLessThanPath?: string;
    NumericLessThanEquals?: number;
    NumericLessThanEqualsPath?: string;
    StringEquals?: string;
    StringEqualsPath?: string;
    StringGreaterThan?: string;
    StringGreaterThanPath?: string;
    StringGreaterThanEquals?: string;
    StringGreaterThanEqualsPath?: string;
    StringLessThan?: string;
    StringLessThanPath?: string;
    StringLessThanEquals?: string;
    StringLessThanEqualsPath?: string;
    TimestampEquals?: string;
    TimestampEqualsPath?: string;
    TimestampGreaterThan?: string;
    TimestampGreaterThanPath?: string;
    TimestampGreaterThanEquals?: string;
    TimestampGreaterThanEqualsPath?: string;
    TimestampLessThan?: string;
    TimestampLessThanPath?: string;
    TimestampLessThanEquals?: string;
    TimestampLessThanEqualsPath?: string;
}

/**
 * Step Functions Choice state. Used to provide branching logic.
 *
 * https://states-language.net/spec.html#choice-state
 */
export interface ChoiceStateNode extends StateNode {
    Type: StatesType.Choice;
    InputPath?: string;
    OutputPath?: string;
    Choices: Array<ChoiceRule>;
    Default?: StateNode;
}

/**
 * Step Functions Wait state. Used to provide blocking logic.
 *
 * https://states-language.net/spec.html#wait-state
 */
export interface WaitStateNode extends StateNode {
    Type: StatesType.Wait;
    InputPath?: string;
    OutputPath?: string;
    Next?: string;
    End?: boolean;
    Seconds?: number;
    Timestamp?: string;
    SecondsPath?: string;
    TimestampPath?: string;
}

/**
 * Step Functions Succeed state. This state is a terminal state which indicites that a state machine has exited successfully.
 * Executions will always end here.
 *
 * https://states-language.net/spec.html#succeed-state
 */
export interface SucceedStateNode extends StateNode {
    Type: StatesType.Succeed;
}

/**
 * Step Functions Fail state. This state is a terminal state which indicites that a state machine execution has failed.
 * Executions will always end here.
 *
 * https://states-language.net/spec.html#fail-state
 */
export interface FailStateNode extends StateNode {
    Type: StatesType.Fail;
    Cause?: string;
    Error?: string;
}

/**
 * An type which defines an executable Step Functions state machine.
 *
 * https://states-language.net/spec.html#toplevelfields
 */
export interface StateMachine extends StateNode {
    Comment?: string;
    StartAt: string;
    States: Record<string, StateNode>;
    Version?: string;
    TimeoutSeconds?: string;
}

/**
 * Step Functions Map state. Runs a set of steps for each element of an input array.
 *
 * https://states-language.net/spec.html#map-state
 */
export interface MapStateNode extends StateNode {
    Type: StatesType.Map;
    Next?: string;
    End?: boolean;
    Iterator: StateMachine;
    ItemsPath?: string;
    MaxConcurrency?: number;
    ResultPath?: string;
    ResultSelector?: string | SerializableRecord;
    Retry?: Array<Retrier>;
    Catch?: Array<Catcher>;
}

/**
 * Step Functions Parallel state. Runs separate state machines in parallel, resolving when all branches are complete.
 *
 * https://states-language.net/spec.html#parallel-state
 */
export interface ParallelStateNode extends StateNode {
    Type: StatesType.Parallel;
    Next?: string;
    End?: boolean;
    Branches: Array<StateMachine>;
    MaxConcurrency?: number;
    ResultPath?: string;
    ResultSelector?: string | SerializableRecord;
    Retry?: Array<Retrier>;
    Catch?: Array<Catcher>;
}
