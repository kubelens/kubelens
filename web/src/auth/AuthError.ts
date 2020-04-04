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
  | 'INTERNAL'
  | 'inactive_session';

export class AuthError extends Error {

  errorCode: AuthErrorCode;

  constructor(errorType: AuthErrorCode, message: string) {
      super();

      this.message = message;
      this.errorCode = errorType;
  }
}
