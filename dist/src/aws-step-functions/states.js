"use strict";
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (Object.hasOwnProperty.call(mod, k)) result[k] = mod[k];
    result["default"] = mod;
    return result;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const Resources = __importStar(require("./resources"));
const config_1 = __importDefault(require("../../config"));
const aws_step_functions_1 = require("../../lib/aws-step-functions");
const { logger } = config_1.default;
/**
* The most generic possible definition of a State node.
*
* Encompasses possibly-liminal States and always-terminal States (the base class).
*/
class State {
    constructor(name) {
        this.name = name;
    }
}
exports.State = State;
/**
* Defines a State node which may, but not necessarily, have downstream nodes.
*/
class LiminalState extends State {
    constructor(name) {
        super(name);
        this._isTerminal = false;
        this._downstreamNode = undefined;
    }
    terminal() {
        this._isTerminal = true;
    }
    setDownstream(state) {
        if (this._isTerminal) {
            throw new Error(`Cannot call setDownstream on terminal state ${this.name}`);
        }
        this._downstreamNode = state;
    }
    _getNextOrEndClause() {
        if (this._isTerminal) {
            return { End: true };
        }
        if (!this._downstreamNode) {
            throw new Error(`Must call setDownstream on non-terminal state ${this.name}`);
        }
        return {
            Next: this._downstreamNode.name
        };
    }
}
class Pass extends LiminalState {
    constructor(name, opts) {
        super(name);
        logger.trace({ name, opts }, 'Created new Pass state');
        this._outputPathPrefix = `$.data.${name}`;
        this._parameters = opts.parameters;
        this._resource = Resources.pass();
        this._result = opts.result;
    }
    index(context) {
        context.registerState({
            name: this.name,
            getOutput: (result) => {
                const resultPath = `${this._outputPathPrefix}.${this._resource.resultPathIdentifier}.result.${result}`;
                return `${resultPath}.${result}`;
            }
        });
    }
    serialize() {
        const passStateNode = Object.assign({ Type: aws_step_functions_1.StatesType.Pass }, this._getNextOrEndClause());
        if (this._parameters) {
            passStateNode.Parameters = this._parameters;
        }
        if (this._result) {
            const resultPath = `${this._outputPathPrefix}.${this._resource.resultPathIdentifier}.result`;
            passStateNode.Result = this._result;
            passStateNode.ResultPath = resultPath;
        }
        return passStateNode;
    }
}
exports.Pass = Pass;
