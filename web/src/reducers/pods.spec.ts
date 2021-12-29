import { reducer, INITIAL_STATE } from './pods';
import { PodActionTypes } from '../actions/pods';

const pod = {
  name: 'name',
  linkedName: 'name',
  namespace: 'namespace',
  pod: {}
};

describe('pods reducer', () => {
  it('should return the initial state', () => {
    expect(reducer(INITIAL_STATE, {
      type: undefined
    })).toEqual(INITIAL_STATE)
  })

  it('should handle GET_POD with default initial state', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: PodActionTypes.GET_POD,
        podOverview: pod
      })
    ).toEqual({
      podOverview: pod
    })
  })

  it('should handle CLEAR_POD', () => {
    expect(
      reducer({
        podOverview: pod
      }, {
        type: PodActionTypes.CLEAR_POD,
        podOverview: pod
      })
    ).toEqual({
      podOverview: pod
    })
  })

  it('should handle SET_SELECTED_CONTAINER_NAME', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: PodActionTypes.SET_SELECTED_CONTAINER_NAME,
        selectedContainerName: 'container name'
      })
    ).toEqual({
      selectedContainerName: 'container name'
    })
  })

  it('should handle SET_SELECTED_POD_NAME', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: PodActionTypes.SET_SELECTED_POD_NAME,
        selectedPodName: 'pod name'
      })
    ).toEqual({
      selectedPodName: 'pod name'
    })
  })
})