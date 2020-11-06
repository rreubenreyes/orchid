// import Emitter from 'events';
import type * as lib from '.';

interface NodeState {
    apply: () => void | Promise<void>;
}

export function createWaitState(director: lib.Directable, data: lib.StateDataWait): NodeState {
    let timeoutInMs: number;

    switch (data.unit) {
        case 'milliseconds': // fallthrough
        case 'millisecond': // fallthrough
        case 'ms':
            timeoutInMs = data.waitTime;
            break;
        case 'seconds': // fallthrough
        case 'second': // fallthrough
        case 'secs': // fallthrough
        case 'sec': // fallthrough
        case 's':
            timeoutInMs = data.waitTime * 1000;
            break;
        case 'minutes': // fallthrough
        case 'minute': // fallthrough
        case 'mins': // fallthrough
        case 'min':
            timeoutInMs = data.waitTime * 1000 * 60;
            break;
        case 'hours': // fallthrough
        case 'hour': // fallthrough
        case 'hrs': // fallthrough
        case 'hr':
            timeoutInMs = data.waitTime * 1000 * 60 * 60;
            break;
        case 'days': // fallthrough
        case 'day': // fallthrough
            timeoutInMs = data.waitTime * 1000 * 60 * 60 * 24;
            break;
    }

    return {
        async apply(): Promise<void> {
            await new Promise((resolve) => {
                setTimeout(resolve, timeoutInMs)
            }).then(() => director.traverse(data.next))
        }
    }
}

export function createTerminusState(director: lib.Directable, data: lib.StateDataTerminus): NodeState {
    return {
        apply(): void {
            director.terminate(data.result === 'success')
        }
    }
}
