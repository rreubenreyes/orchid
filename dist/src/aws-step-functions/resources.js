"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
function pass() {
    return {
        resourceName: 'pass',
        resultPathIdentifier: 'pass',
    };
}
exports.pass = pass;
exports.lambda = {
    lnvoke: ({ functionName }) => ({
        resourceName: 'arn:aws:states:::lambda:invoke',
        resultPathIdentifier: 'lambdaResult',
        parameters: {
            FunctionName: functionName,
        }
    }),
    lnvokeWaitForTaskToken: ({ functionName }) => ({
        resourceName: 'arn:aws:states:::lambda:invoke.waitForTaskToken',
        resultPathIdentifier: 'lambdaResult',
        parameters: {
            FunctionName: functionName,
        }
    }),
};
