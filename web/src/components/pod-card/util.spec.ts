import { PodPhaseStyle, PENDING, RUNNING, SUCCEEDED, FAILED, UNKNOWN } from './util';

describe('overview util should', () => {
  test('return pending', () => {
    const phase = PodPhaseStyle(PENDING);
    expect(phase.class).toBe('pod-container pod-pending')
  })

  test('return pending', () => {
    const phase = PodPhaseStyle(RUNNING);
    expect(phase.class).toBe('pod-container pod-running')
  })

  test('return succeeded', () => {
    const phase = PodPhaseStyle(SUCCEEDED);
    expect(phase.class).toBe('pod-container pod-succeeded')
  })

  test('return failed', () => {
    const phase = PodPhaseStyle(FAILED);
    expect(phase.class).toBe('pod-container pod-failed')
  })

  test('return unknown', () => {
    const phase = PodPhaseStyle(UNKNOWN);
    expect(phase.class).toBe('pod-container pod-unknown')
  })

  test('return pending', () => {
    const phase = PodPhaseStyle('none');
    expect(phase.class).toBeUndefined()
  })
})