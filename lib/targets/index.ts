type StateKind = 'pass'
| 'wait'
| 'choice'
| 'effect'
// | 'iterator'
// | 'parallel'
| 'terminus'


interface StateData {
    kind: StateKind;
    next: StateData;
}

export interface Directable {
    traverse: (next: StateData) => void;
    terminate: (successful: boolean) => void;
    terminated: boolean;
    successful: boolean;
}

export type StateDataPass = StateData
export interface StateDataWait extends StateData {
    kind: 'wait';
    waitTime: number;
    unit: 'ms'
    | 'milliseconds'
    | 'millisecond'
    | 'seconds'
    | 'second'
    | 'secs'
    | 'sec'
    | 's'
    | 'minutes'
    | 'minute'
    | 'mins'
    | 'min'
    | 'hours'
    | 'hour'
    | 'hrs'
    | 'hr'
    | 'days'
    | 'day';
}

export interface Choice {
    name: string;
    variable: unknown;
    $eq?: unknown;
    $lt?: unknown;
    $lte?: unknown;
    $gt?: unknown;
    $gte?: unknown;
    $between?: [unknown, unknown];
}

export interface StateDataChoice extends Pick<StateData, 'kind'> {
    choices: Choice[];
    predicate: {
        $and: Array<string | '$or' | '$and'>;
        $nand: Array<string | '$or' | '$and'>;
        $or: Array<string | '$or' | '$and'>;
    };
}

export interface StateDataTerminus extends Pick<StateData, 'kind'> {
    result: 'success' | 'failure';
}

export interface StateDataEffect extends StateData {
    effect: unknown;
}
