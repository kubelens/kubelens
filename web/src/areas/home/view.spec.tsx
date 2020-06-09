import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import HomePage from './view'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    match: {
      params: {
        appName: 'appname'
      }
    },
    onFilterChanged: jest.fn(),
    onViewApp: jest.fn(),
    filteredApps: [
      {name: 'name1', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'},
      {name: 'name2', namespace: 'namespace', labelSelector: 'labelSelector', kind: 'kind'}
    ],
    selectedAppName: 'appname',
    isLoading: false
  }

  const wrapper = shallow(<HomePage {...props} />)

  return {
    props,
    wrapper
  }
}

describe('home view should', () => {
  test('render with filtered apps', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Input').length).toBe(1);
    expect(wrapper.find('ReactList').length).toBe(1);
    expect(wrapper.find('Route').length).toBe(3);
  })
})