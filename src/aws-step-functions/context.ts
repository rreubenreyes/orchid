export interface OrchidIndex {
    name: string;
    getOutput(result: string): string;
}

export interface OrchidContext {
    registerState(stateIndex: OrchidIndex): void;
    getState(stateName: string): OrchidIndex;
}

export interface IndexedState {
    serialize<T>(): T;
}

