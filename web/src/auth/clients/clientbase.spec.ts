import * as Exports from './clientbase';

describe('ClientBase', () => {
  let cb: Exports.ClientBase;
  const testToken = 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiIsImtpZCI6Ik56SkdNVGcyUWpWRE1EQXlPRVJHUmpjeE1UWkNNa1pET1VZME9VVTNSamswUVRreFJEY3dOZyJ9.eyJpc3MiOiJodHRwczovL2F1dGguY29tcGFueS5jb20vIiwic3ViIjoiOW1EbHo5NGFMb2luOVgzeUZYWEVGbXNtRVhyMmVzaDRAY2xpZW50cyIsImF1ZCI6Imh0dHA6Ly9hcGkuY29tcGFueS5jb20iLCJpYXQiOjE1MzAxMDcyMTQsImV4cCI6MTUzMDEwNzIxNCwiZ3R5IjoiY2xpZW50LWNyZWRlbnRpYWxzIn0.1Hf3vMaoVHaQxM4wIWWmDVg144oxjxkDpcMoxWJZOkM';
  const noExpirationClaimTestToken = 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiIsImtpZCI6Ik56SkdNVGcyUWpWRE1EQXlPRVJHUmpjeE1UWkNNa1pET1VZME9VVTNSamswUVRreFJEY3dOZyJ9.eyJpc3MiOiJodHRwczovL2F1dGguY29tcGFueS5jb20vIiwic3ViIjoiOW1EbHo5NGFMb2luOVgzeUZYWEVGbXNtRVhyMmVzaDRAY2xpZW50cyIsImF1ZCI6Imh0dHA6Ly9hcGkuY29tcGFueS5jb20iLCJpYXQiOjE1MzAxMDcyMTQsImd0eSI6ImNsaWVudC1jcmVkZW50aWFscyJ9.UtkKmiPB8BUBHJHRFOAYFIOCeK3bOZsb-aWOv_AMF58';

  beforeEach(() => {
    cb = new Exports.ClientBase();
  });

  it('should return expiration date from specified token', () => {
    const expectedExpirationDate = new Date('2018-06-27T13:46:54.000Z');
    const expirationDateResult = cb.getTokenExpirationDate(testToken);
    expect(expirationDateResult).toEqual(expectedExpirationDate);
  });

  it('should indicate if the token has expired', () => {
    const isTokenExpiredResult = cb.isTokenExpired(testToken);
    expect(isTokenExpiredResult).toBeTruthy();
  });

  it('should return null if the specified token does not contains an expiration claim', () => {
    const expirationDateResult = cb.getTokenExpirationDate(noExpirationClaimTestToken);
    expect(expirationDateResult).toBeNull();
  });
});
