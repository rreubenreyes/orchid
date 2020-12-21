import _ from 'lodash';

import * as Resources from './resources';
import config from '../config';
import { StatesType } from '../lib/aws-step-functions';
import type * as SFN from '../lib/aws-step-functions';
import type { PassResource } from './resources';
import type { OrchidContext, ReadonlyOrchidContext, IndexedState } from './context';

const { logger } = config;

/**
* The most generic possible definition of a State node.
*
* Encompasses possibly-liminal States and always-terminal States (the base class).
*/
export abstract class State {
    public name: string;

    constructor(name: string) {
        this.name = name;
    }

    abstract index(context: OrchidContext): IndexedState<SFN.StateNode>;
}

/**
* Defines a State node which may, but not necessarily, have downstream nodes.
*/
abstract class LiminalState extends State {
    protected _isTerminal = false;
    protected _downstreamNode: State | undefined = undefined;

    constructor(name: string) {
        super(name);
    }

    terminal(): void {
        this._isTerminal = true;
    }

    setDownstream(state: State): void {
        if (this._isTerminal) {
            throw new Error(`Cannot call setDownstream on terminal state ${this.name}`);
        }

        this._downstreamNode = state;
    }

    protected _getNextOrEndClause():
    { Next: string; End?: boolean } | { Next?: string; End: boolean } {
        if (this._isTerminal) {
            return { End: true };
        }

        if (!this._downstreamNode) {
            throw new Error(`Must call setDownstream on non-terminal state ${this.name}`)
        }

        return {
            Next: this._downstreamNode.name
        }
    }
}

export class Pass extends LiminalState {
    private _outputPathPrefix: string;
    private _result?: SFN.Serializable;
    private _task?: PassResource;

    constructor(name: string, opts: {
        task?: PassResource;
    } = {}) {
        super(name);

        logger.trace({ name, opts }, 'Created new Pass state');

        this._outputPathPrefix = `$.data.${name}`;
        this._task = opts.task;
    }

    index(context: OrchidContext): IndexedState<SFN.PassStateNode> {
        context.registerState({
            name: this.name,
            getOutput: (result: string) => {
                if (!this._task.result) {
                    throw new Error(`State ${this.name} does not produce an output`);
                }

                const resultPath = `${this._outputPathPrefix}.${this._resource.resultPathIdentifier}.result`;

                return `${resultPath}.${result}`
            }
        });

        return {
            serialize: (): SFN.PassStateNode => {
                const passStateNode: SFN.PassStateNode = {
                    Type: StatesType.Pass,
                    ...this._getNextOrEndClause(),
                }

                if (this._parameters) {
                    const readonlyContext = _.pick(context, 'getState');
                    passStateNode.Parameters = this._parameters(readonlyContext);
                }

                if (this._result) {
                    const resultPath = `${this._outputPathPrefix}.${this._resource.resultPathIdentifier}.result`;

                    passStateNode.Result = this._result;
                    passStateNode.ResultPath = resultPath;
                }

                return passStateNode;
            }
        }
    }
}

