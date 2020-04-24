import { reducer, INITIAL_STATE } from './error';
import { ErrorActionTypes } from '../actions/error';

describe('error reducer', () => {
  it('should return the initial state', () => {
    expect(reducer(INITIAL_STATE, {
      type: undefined,
      apiOpen: false,
      status: 123,
      statusText: 'should be seen',
      message: 'no message'
    })).toEqual(INITIAL_STATE)
  })

  it('should handle OPEN_API_ERROR_MODAL with default initial state', () => {
    expect(
      reducer({
        apiOpen: false
      }, {
        type: ErrorActionTypes.OPEN_API_ERROR_MODAL,
        apiOpen: true,
        status: 401,
        statusText: 'unauthorized',
        message: 'message'
      })
    ).toEqual({
      apiOpen: true,
      status: 401,
      statusText: 'unauthorized',
      message: 'message'
    })
  })

  it('should handle CLOSE_API_ERROR_MODAL with default initial state', () => {
    expect(
      reducer({
        apiOpen: true
      }, {
        type: ErrorActionTypes.CLOSE_API_ERROR_MODAL,
        apiOpen: false
      })
    ).toEqual({
      apiOpen: false
    })
  })
})