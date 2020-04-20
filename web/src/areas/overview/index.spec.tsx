import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import { Overview } from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    identityToken: 'token',
    selectedAppName: 'appname',
    getAppOverview: jest.fn(),
    setSelectedAppName: jest.fn(),
    overviewsEmpty: false,
    error: {
      apiOpen: false
    },
    isLoading: false,
    match: {
      params: {
        appName: 'appname'
      }
    },
    location: {
      search: 'query'
    }
  }

  const wrapper = shallow(<Overview {...props} />)

  return {
    props,
    wrapper
  }
}

describe('overview should', () => {
  test('render', () => {
    const { wrapper } = setup();

    expect(wrapper.find('div.overview-container').length).toBe(1);
    expect(wrapper.find('APIErrorModal').length).toBe(1);
  })
})