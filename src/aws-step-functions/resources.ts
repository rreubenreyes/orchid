export interface IntermediateResource {
    resourceName: string;
    resultPathIdentifier: string;
    parameters?: Record<string, string>;
}

export function pass(): IntermediateResource {
    return {
        resourceName: 'pass',
        resultPathIdentifier: 'pass',
    }
}

export const lambda = {
    invoke: ({ functionName }: { functionName: string }): IntermediateResource => ({
        resourceName: 'arn:aws:states:::lambda:invoke',
        resultPathIdentifier: 'lambdaResult',
        parameters: {
            FunctionName: functionName,
        }
    }),
    invokeWaitForTaskToken: ({ functionName }: { functionName: string }): IntermediateResource => ({
        resourceName: 'arn:aws:states:::lambda:invoke.waitForTaskToken',
        resultPathIdentifier: 'lambdaResult',
        parameters: {
            FunctionName: functionName,
        }
    }),
}


