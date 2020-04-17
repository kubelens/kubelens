import { reducer, INITIAL_STATE } from './pods';
import { PodActionTypes } from '../actions/pods';

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
        pod: {
          name: "name",
          namespace: "namespace",
          hostIP: "hostip",
          podIP: "podip",
          startTime: "starttime",
          phase: {},
          phaseMessage: 'phasemessage',
          containerStatus: [],
          status: {
            containerStatuses: [],
            conditions: [],
            startTime: "",
            phase: ""
          },
          spec: {
            containers: [{
              env: {}
            }]
          },
          containerNames: []
        }
      })
    ).toEqual({
      pod: {
        name: "name",
        namespace: "namespace",
        hostIP: "hostip",
        podIP: "podip",
        startTime: "starttime",
        phase: {},
        phaseMessage: 'phasemessage',
        containerStatus: [],
        status: {
          containerStatuses: [],
          conditions: [],
          startTime: "",
          phase: ""
        },
        spec: {
          containers: [{
            env: {}
          }]
        },
        containerNames: []
      }
    })
  })

  it('should handle CLEAR_POD', () => {
    expect(
      reducer({
        pod: {
          name: "name",
          namespace: "namespace",
          hostIP: "hostip",
          podIP: "podip",
          startTime: "starttime",
          phase: {},
          phaseMessage: 'phasemessage',
          containerStatus: [],
          status: {
            containerStatuses: [],
            conditions: [],
            startTime: "",
            phase: ""
          },
          spec: {
            containers: [{
              env: {}
            }]
          },
          containerNames: []
        }
      }, {
        type: PodActionTypes.CLEAR_POD,
        pod: null
      })
    ).toEqual({
      pod: null
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