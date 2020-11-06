class StateNode {
    public name: string;

    constructor({ name }: { name: string }) {
        this.name = name;
    }
}

export default class StateMachine {
    public static targets = CompileTargets
}
