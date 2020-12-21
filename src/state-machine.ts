import { createContext } from './context';

import type { State } from './states';
import type { IndexedState } from './context';
import type * as SFN from '../lib/aws-step-functions';

interface StateMachine {
    definition: SFN.Serializable;
    // asJson(): string;
}

export function create({ startAt, states }: {
    startAt: State | string;
    states: Array<IndexedState<SFN.StateNode>>;
}): StateMachine {
    const context = createContext();

    return {
        definition: {},
        // asJson: JSON.stringify,
    }
}
