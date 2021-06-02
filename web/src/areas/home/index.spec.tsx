import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import { Home } from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    cluster: 'minikube',
    identityToken: 'token',
    selectedOverview: 'appname',
    filteredApps: [{name: 'appname'}],
    apps: [{name: 'appname'}],
    getApps: jest.fn(),
    getAppOverview: jest.fn(),
    setselectedOverview: jest.fn(),
    filterApps: jest.fn(),
    error: {
      apiOpen: false
    },
    isLoading: false,
    match: {
      params: {
        appName: 'appname'
      }
    }
  }

  const wrapper = shallow(<Home {...props} />)

  return {
    props,
    wrapper
  }
}

describe('home should', () => {
  test('render', () => {
    const { wrapper } = setup();

    expect(wrapper.find('APIErrorModal').length).toBe(1);
  })
})