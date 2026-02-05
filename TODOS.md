# TODO Tracking

This document tracks all TODO items found in the codebase. These should be converted into GitHub issues for proper tracking.

## Security Handler TODOs

### TODO 1: Validate Bearer Token
- **File**: `authn/service/security-handler.go:26`
- **Context**: HandleBearerAuth method
- **Description**: Validate the token and possibly add user info to the context
- **Priority**: High (Security-related)
- **Current Behavior**: Currently logs the token but doesn't validate it
- **Code Reference**:
```go
func (s SecurityService) HandleBearerAuth(
	ctx context.Context,
	operationName api.OperationName,
	t api.BearerAuth,
) (context.Context, error) {
	slog.Info("Bearer auth received", "operation", operationName, "token", t.Token)
	// TODO:  validate the token and possibly add user info to the context.
	return ctx, nil
}
```

## Authentication Service TODOs

### TODO 2: Implement GetUserProfile
- **File**: `authn/service/auth-handler.go:14-16`
- **Context**: GetUserProfile method
- **Description**: Implement user profile retrieval functionality
- **Priority**: High
- **Current Behavior**: Panics with "implement me"

### TODO 3: Implement LoginUser
- **File**: `authn/service/auth-handler.go:19-22`
- **Context**: LoginUser method
- **Description**: Implement user login functionality
- **Priority**: High
- **Current Behavior**: Panics with "implement me"

### TODO 4: Implement LogoutUser
- **File**: `authn/service/auth-handler.go:24-27`
- **Context**: LogoutUser method
- **Description**: Implement user logout functionality
- **Priority**: High
- **Current Behavior**: Panics with "implement me"

### TODO 5: Implement RegisterUser
- **File**: `authn/service/auth-handler.go:29-32`
- **Context**: RegisterUser method
- **Description**: Implement user registration functionality
- **Priority**: High
- **Current Behavior**: Panics with "implement me"

### TODO 6: Implement RequestPasswordReset
- **File**: `authn/service/auth-handler.go:34-37`
- **Context**: RequestPasswordReset method
- **Description**: Implement password reset request functionality
- **Priority**: Medium
- **Current Behavior**: Panics with "implement me"

### TODO 7: Implement ResetPassword
- **File**: `authn/service/auth-handler.go:39-42`
- **Context**: ResetPassword method
- **Description**: Implement password reset confirmation functionality
- **Priority**: Medium
- **Current Behavior**: Panics with "implement me"

### TODO 8: Implement UpdateUserProfile
- **File**: `authn/service/auth-handler.go:44-47`
- **Context**: UpdateUserProfile method
- **Description**: Implement user profile update functionality
- **Priority**: Medium
- **Current Behavior**: Panics with "implement me"

## Notes

- All authentication service methods are currently unimplemented stubs that panic
- The security handler accepts any bearer token without validation
- These TODOs should be converted to GitHub issues and assigned to team members
- Consider breaking down larger TODOs into smaller, manageable issues
