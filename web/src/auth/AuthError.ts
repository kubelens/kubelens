/*
MIT License

Copyright (c) 2019 The KubeLens Authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
export type AuthErrorCode = 'access_denied'
  | 'missing_configuration'
  | 'login_required'
  | 'invalid_request'
  | 'unauthorized_client'
  | 'server_error'
  | 'invalid_scope'
  | 'temporarily_unavailable'
  | 'unsupported_response_type'
  | 'too_many_attempts'
  | 'invalid_user_password'
  | 'INTERNAL';



export class AuthError extends Error {

  errorCode: AuthErrorCode;

  constructor(errorType: AuthErrorCode, message: string) {
    super();

    this.message = message;
    this.errorCode = errorType;
  }
}