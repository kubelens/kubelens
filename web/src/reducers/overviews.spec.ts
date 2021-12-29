import { reducer, INITIAL_STATE } from './overviews';
import { OverviewActionTypes } from '../actions/overviews';

const ov = {
  linkedName: 'name',
  namespace: 'namespace',
  pods: [],
  services: [],
  daemonSets: [],
  jobs: [],
  deployments: [],
  configMaps: [],
  replicaSets: []
};

const sn = {
  linkedName: 'name',
  namespace: 'namespace'
};

describe('overviews reducer', () => {
  it('should return the initial state', () => {
    expect(reducer(INITIAL_STATE, {type: undefined})).toEqual({})
  })

  it('should handle GET_OVERVIEW with default initial state', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: OverviewActionTypes.GET_OVERVIEW,
        overview: ov
      })
    ).toEqual({
      overview: ov
    })
  })

  it('should handle GET_OVERVIEWS with initial state poplulated', () => {
    expect(
      reducer({
        overviews: [ov]
      }, {
        type: OverviewActionTypes.GET_OVERVIEWS,
        overviews: [ov]
      })
    ).toEqual({
      filteredOverviews: [ov],
      overviews: [ov]
    })
  })

  it('should handle GET_OVERVIEWS with default initial state', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: OverviewActionTypes.GET_OVERVIEWS,
        overviews: [ov]
      })
    ).toEqual({
      overviews: [ov],
      filteredOverviews: [ov]
    })
  })

  it('should handle FILTER_OVERVIEWS with default initial state', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: OverviewActionTypes.FILTER_OVERVIEWS,
        filteredOverviews: [ov]
      })
    ).toEqual({
      filteredOverviews: [ov]
    })
  })

  it('should handle FILTER_OVERVIEWS with initial state poplulated', () => {
    expect(
      reducer({
        selectedOverview: sn
      }, {
        type: OverviewActionTypes.FILTER_OVERVIEWS,
        filteredOverviews: [ov]
      })
    ).toEqual({
      selectedOverview: sn,
      filteredOverviews: [ov]
    })
  })
})