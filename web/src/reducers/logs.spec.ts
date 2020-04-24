import { reducer, INITIAL_STATE } from './logs';
import { LogsActionTypes } from '../actions/logs';

describe('logs reducer', () => {
  it('should return the initial state', () => {
    expect(reducer(INITIAL_STATE, {
      type: undefined
    })).toEqual(INITIAL_STATE)
  })

  it('should handle GET_LOGS with default initial state', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: LogsActionTypes.GET_LOGS,
        logs: {
          pod: 'podname',
          output: 'log output'
        }
      })
    ).toEqual({
      logStreamEnabled: false,
      logsError: undefined,
      logs: {
        pod: 'podname',
        output: 'log output'
      }
    })
  })

  it('should handle TOGGLE_LOG_STREAM true', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: LogsActionTypes.TOGGLE_LOG_STREAM,
        logStreamEnabled: true
      })
    ).toEqual({
      logStreamEnabled: true
    })
  })

  it('should handle TOGGLE_LOG_STREAM false', () => {
    expect(
      reducer({
        logStreamEnabled: true
      }, {
        type: LogsActionTypes.TOGGLE_LOG_STREAM,
        logStreamEnabled: false
      })
    ).toEqual({
      logStreamEnabled: false
    })
  })

  it('should handle CLEAR_ERRORS', () => {
    expect(
      reducer({
        logStreamEnabled: false,
        logsError: new Error('test error')
      }, {
        type: LogsActionTypes.CLEAR_ERRORS,
        logsError: null
      })
    ).toEqual({
      logStreamEnabled: false,
      logsError: undefined
    })
  })
})