"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createTerminusState = exports.createWaitState = void 0;
function createWaitState(director, data) {
    let timeoutInMs;
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
        apply() {
            return __awaiter(this, void 0, void 0, function* () {
                yield new Promise((resolve) => {
                    setTimeout(resolve, timeoutInMs);
                }).then(() => director.traverse(data.next));
            });
        }
    };
}
exports.createWaitState = createWaitState;
function createTerminusState(director, data) {
    return {
        apply() {
            director.terminate(data.result === 'success');
        }
    };
}
exports.createTerminusState = createTerminusState;
