import type { ReadonlyOrchidContext } from './context';
import type * as SFN from '../lib/aws-step-functions';

type ParametersGetter = (context: ReadonlyOrchidContext) => SFN.Serializable;

export interface IntermediateResource {
    resourceName: string;
    resultPathIdentifier: string;
    parameters?: ParametersGetter;
}

export interface PassResource extends IntermediateResource {
    result?: SFN.Serializable;
};

export function pass({ parameters, result }: {
    parameters?: ParametersGetter;
    result?: SFN.Serializable;
}): PassResource {
    return {
        resourceName: 'pass',
        resultPathIdentifier: 'pass',
        parameters,
        result,
    }
}

// export const lambda = {
//     invoke: ({ functionName }: { functionName: string }): IntermediateResource => ({
//         resourceName: 'arn:aws:states:::lambda:invoke',
//         resultPathIdentifier: 'lambdaResult',
//         parameters: {
//             FunctionName: functionName,
//         }
//     }),
//     invokeWaitForTaskToken: ({ functionName }: { functionName: string }): IntermediateResource => ({
//         resourceName: 'arn:aws:states:::lambda:invoke.waitForTaskToken',
//         resultPathIdentifier: 'lambdaResult',
//         parameters: {
//             FunctionName: functionName,
//         }
//     }),
// }


