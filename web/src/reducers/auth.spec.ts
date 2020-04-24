import { reducer, INITIAL_STATE } from './auth';
import { AuthActionTypes } from '../actions/auth';

describe('auth reducer', () => {
  it('should return the initial state', () => {
    expect(reducer(INITIAL_STATE, {
      type: undefined,
      identityToken: 'invalid'
    })).toEqual(INITIAL_STATE)
  })

  it('should handle SET_IDENTITY_TOKEN with default initial state', () => {
    expect(
      reducer({
        identityToken: 'old token'
      }, {
        type: AuthActionTypes.SET_IDENTITY_TOKEN,
        identityToken: 'new token'
      })
    ).toEqual({
      identityToken: 'new token'
    })
  })
})