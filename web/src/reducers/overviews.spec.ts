import { reducer, INITIAL_STATE } from './apps';
import { AppsActionTypes } from '../actions/apps';

describe('apps reducer', () => {
  it('should return the initial state', () => {
    expect(reducer(INITIAL_STATE, {type: undefined})).toEqual({})
  })

  it('should handle GET_APP_OVERVIEW with default initial state', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: AppsActionTypes.GET_APP_OVERVIEW,
        appOverview: [{
          podOverviews: [],
          serviceOvervies: [],
          daemonSetOverviews: []
        }]
      })
    ).toEqual({
      appOverview: [{
        podOverviews: [],
        serviceOvervies: [],
        daemonSetOverviews: []
      }]
    })
  })

  it('should handle GET_APP_OVERVIEW with initial state poplulated', () => {
    expect(
      reducer({
        selectedAppName: 'appname',
        filteredApps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}],
        apps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}]
      }, {
        type: AppsActionTypes.GET_APP_OVERVIEW,
        appOverview: [{
          podOverviews: [],
          serviceOvervies: [],
          daemonSetOverviews: []
        }]
      })
    ).toEqual({
      selectedAppName: 'appname',
      filteredApps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}],
      apps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}],
      appOverview: [{
        podOverviews: [],
        serviceOvervies: [],
        daemonSetOverviews: []
      }]
    })
  })

  it('should handle GET_APPS with default initial state', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: AppsActionTypes.GET_APPS,
        apps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}]
      })
    ).toEqual({
      apps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}],
      filteredApps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}]
    })
  })

  it('should handle GET_APPS with initial state poplulated', () => {
    expect(
      reducer({
        selectedAppName: 'appname',
        filteredApps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}]
      }, {
        type: AppsActionTypes.GET_APPS,
        apps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}]
      })
    ).toEqual({
      selectedAppName: 'appname',
      filteredApps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}],
      apps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}]
    })
  })

  it('should handle FILTER_APPS with default initial state', () => {
    expect(
      reducer(INITIAL_STATE, {
        type: AppsActionTypes.FILTER_APPS,
        filteredApps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}]
      })
    ).toEqual({
      filteredApps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}]
    })
  })

  it('should handle FILTER_APPS with initial state poplulated', () => {
    expect(
      reducer({
        selectedAppName: 'appname'
      }, {
        type: AppsActionTypes.FILTER_APPS,
        filteredApps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}]
      })
    ).toEqual({
      selectedAppName: 'appname',
      filteredApps: [{name: 'name', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}]
    })
  })
})