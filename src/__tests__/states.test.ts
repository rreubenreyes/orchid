import * as States from '../states';
import { createContext } from '../context';

import { StatesType } from '../../lib/aws-step-functions';

describe('States', () => {
    describe('Pass', () => {
        it('can be instantiated', () => {
            const passState = new States.Pass('TestPassState');
            expect(passState.name).toBe('TestPassState');
        })

        it('compiles down to SFN.PassStateNode', () => {
            const context = createContext();
            const passState = new States.Pass('TestPassState');

            passState.terminal();

            const serialized = passState
                .index(context)
                .serialize();

            console.log(serialized);

            expect(serialized.Type).toBe(StatesType.Pass);
        })

        it('serializes `result` parameter to PassStateNode.Result', () => {
            const context = createContext();
            const testResultClause = { works: true }
            const stateName = 'TestPassState'
            const passState = new States.Pass(stateName, {
                result: testResultClause
            });

            passState.terminal();

            const serialized = passState
                .index(context)
                .serialize();

            console.log(serialized);

            expect(serialized.Result).toEqual(testResultClause);
            expect(serialized.ResultPath).toContain(stateName);
        })

        it('accepts and correctly serializes parameters', () => {
            const context = createContext();
            const testResultClause = { works: true }

            const firstState = new States.Pass('FirstState', {
                result: testResultClause
            });

            const secondState = new States.Pass('SecondState', {
                parameters: (context) => ({
                    testParameter: context
                        .getState('FirstState')
                        .getOutput('works')
                })
            })

            firstState.setDownstream(secondState);
            secondState.terminal();

            firstState.index(context);

            const serialized = secondState.index(context).serialize();

            console.log(serialized);

            expect(serialized).toBeTruthy();
        })
    })
})
