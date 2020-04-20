import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import { Cluster } from '.'
// import { NavBar } from '../components/nav'
// import { Home } from './home'

Enzyme.configure({ adapter: new Adapter() })

function setup() {
  const props = {
    clustersState: {
      cluster: 'minikube'
    }
  }

  const enzymeWrapper = shallow(<Cluster {...props} />)

  return {
    props,
    enzymeWrapper
  }
}

describe('main area index should', () => {
  test('render Cluster', () => {
    const { enzymeWrapper } = setup();

    expect(enzymeWrapper.find('div.background-logo').length).toBe(1);
    // expect(enzymeWrapper.find(NavBar).length).toBe(1);
    // expect(enzymeWrapper.find(Home).length).toBe(1);
  })

  // it('should call addTodo if length of text is greater than 0', () => {
  //   const { enzymeWrapper, props } = setup()
  //   const input = enzymeWrapper.find('TodoTextInput')
  //   input.props().onSave('')
  //   expect(props.addTodo.mock.calls.length).toBe(0)
  //   input.props().onSave('Use Redux')
  //   expect(props.addTodo.mock.calls.length).toBe(1)
  // })
})