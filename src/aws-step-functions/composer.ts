import _ from 'lodash';

import { StatesType } from '../../lib/aws-step-functions';
import type * as SFN from '../../lib/aws-step-functions';

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

class Pass {
    public name: string;
    private definition: Record<string, unknown>;

    constructor(name: string, definition: {
        parameters?: SFN.Serializable;
        result?: SFN.Serializable;
    }) {
        this.name = name;
        this.definition = definition;
    }

    serialize() {
        const partiallySerializedDefinition = Object.fromEntries(
            Object.entries(this.definition).map([k, v] => ([_.camelCase(k), v]))
        )
    }
}


/**
 *
 * const firstState = new States.Pass('FirstState', {
 *
 * })
 * const secondState = new State().terminal();
 *
 * firstState.setDownstream(secondState)
 */
