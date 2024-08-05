package gatekeeper

import "github.com/gin-gonic/gin"

type DefaultAuthorizer struct {
	roleToPermissionMap map[string]map[string]bool
	accessProvider      AccessProvider
	applicationName     string
}

func NewDefaultAuthorizer() Authorizer {
	//read from config
	return &DefaultAuthorizer{
		roleToPermissionMap: nil,
	}
}
func NewDefaultAuthorizerWithPermissionProvider(accessProvider AccessProvider, options Options, applicationName string) Authorizer {

	return &DefaultAuthorizer{
		roleToPermissionMap: nil,
		accessProvider:      accessProvider,
		applicationName:     applicationName,
	}
}

func (authorizer *DefaultAuthorizer) Authorize(opts ...interface{}) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		roles := ctx.GetStringSlice("Roles")
		for _, role := range roles {
			permissions := authorizer.roleToPermissionMap[role]
			for _, opt := range opts {
				permission, ok := opt.(string)
				if ok && permissions[permission] {
					ctx.Set("AuthContext", authorizer.getAuthContextForAllRoles(roles))
					ctx.Next()
					return
				}
			}
		}
		ctx.Abort()
		return
	}
}

func (authorizer *DefaultAuthorizer) getAuthContextForAllRoles(roles []string) map[string]map[string]bool {
	authContext := map[string]map[string]bool{}
	for _, role := range roles {
		authContext[role] = authorizer.roleToPermissionMap[role]
	}
	return authContext
}

func (authorizer *DefaultAuthorizer) VerifyAuthorization(ctx *gin.Context, opt ...AccessValidator) (bool, error) {

	isAuthorized, err := authorizer.accessProvider.ValidateAccess(ctx)
	if err != nil {
		return isAuthorized, err
	}
	return isAuthorized, nil
}
