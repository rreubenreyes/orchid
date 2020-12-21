export interface OrchidIndexEntry {
    name: string;
    getOutput(result: string): string;
}

export interface OrchidContext {
    readonly index: Record<string, OrchidIndexEntry>;
    registerState(state: OrchidIndexEntry): void;
    getState(stateName: string): OrchidIndexEntry;
}

export type ReadonlyOrchidContext = Pick<OrchidContext, 'getState'>

export interface IndexedState<T> {
    serialize(): T;
}

export const createContext = (): OrchidContext => {
    const _index: Record<string, OrchidIndexEntry> = {};

    return {
        get index(): Record<string, OrchidIndexEntry> {
            return _index;
        },
        registerState(state): void {
            _index[state.name] = state;
        },
        getState(stateName): OrchidIndexEntry {
            return _index[stateName];
        }
    }
}
