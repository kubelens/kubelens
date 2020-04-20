import { reducer, INITIAL_STATE } from './cluster';
import { ClustersActionTypes } from '../actions/cluster';

describe('cluster reducer', () => {
  it('should return the initial state', () => {
    expect(reducer(INITIAL_STATE, {
      type: undefined,
      cluster: 'none'
    })).toEqual(INITIAL_STATE)
  })

  it('should handle SET_SELECTED_CLUSTER with default initial state', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: ClustersActionTypes.SET_SELECTED_CLUSTER,
        cluster: 'clustername'
      })
    ).toEqual({
      cluster: 'clustername'
    })
  })
})