import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import { NavBar } from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = { 
    isLoading: false,
    selectedCluster: 'minikube',
    identityToken: 'token',
    setSelectedCluster: jest.fn(),
    getApps: jest.fn()
  }

  const wrapper = shallow(<NavBar {...props} />)

  return {
    props,
    wrapper
  }
}

describe('NavBar should', () => {
  test('render', () => {
    const { wrapper } = setup();

    expect(wrapper.find('div.kubelens-navbar').length).toBe(1);
  })

  test('have logo image present', () => {
    const { wrapper } = setup();

    expect(wrapper.find('img')).toBeDefined();
    const { alt } = wrapper.find('img').props();
    expect(alt).toBe('Kubelens Logo');
  })

  test('have cluster selection dropdown on right', () => {
    const { wrapper } = setup();

    expect(wrapper.find('#navbar-right').length).toBe(1);
    expect(wrapper.find('#navbar-right').children().find('Dropdown')).toBeDefined();
  })
})