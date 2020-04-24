import { reducer, INITIAL_STATE } from './loading';
import { LoadingActionTypes } from '../actions/loading';

describe('loading reducer', () => {
  it('should return the initial state', () => {
    expect(reducer(INITIAL_STATE, {
      type: undefined, 
      loading: true
    })).toEqual(INITIAL_STATE)
  })

  it('should handle LOADING true', () => {
    expect(
      reducer({
        loading: false
      }, {
        type: LoadingActionTypes.LOADING,
        loading: true
      })
    ).toEqual({
      loading: true
    })
  })

  it('should handle LOADING false', () => {
    expect(
      reducer({
        loading: true
      }, {
        type: LoadingActionTypes.LOADING,
        loading: false
      })
    ).toEqual({
      loading: false
    })
  })
})