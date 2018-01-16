// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/mattermost/mattermost-server/model"
)

func (a *App) GetRole(id string) (*model.Role, *model.AppError) {
	if result := <-a.Srv.Store.Role().Get(id); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.Role), nil
	}
}

func (a *App) GetRoleByName(name string) (*model.Role, *model.AppError) {
	if result := <-a.Srv.Store.Role().GetByName(name); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.Role), nil
	}
}

func (a *App) GetRolesByNames(names []string) ([]*model.Role, *model.AppError) {
	if result := <-a.Srv.Store.Role().GetByNames(names); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.([]*model.Role), nil
	}
}

func (a *App) PatchRole(role *model.Role, patch *model.RolePatch) (*model.Role, *model.AppError) {
	role.Patch(patch)
	role, err := a.UpdateRole(role)
	if err != nil {
		return nil, err
	}

	return role, err
}

func (a *App) UpdateRole(role *model.Role) (*model.Role, *model.AppError) {
	if result := <-a.Srv.Store.Role().Save(role); result.Err != nil {
		return nil, result.Err
	} else {
		// TODO: Is any cache invalidation required here?
		a.sendUpdatedRoleEvent(role)

		return role, nil
	}
}

func (a *App) CheckRolesExist(roleNames []string) (bool, *model.AppError) {
	roles, err1 := a.GetRolesByNames(roleNames)
	if err1 != nil {
		return false, err1
	}

	for _, name := range roleNames {
		nameFound := false
		for _, role := range roles {
			if name == role.Name {
				nameFound = true
				break
			}
		}
		if !nameFound {
			return false, nil
		}
	}

	return true, nil
}

func (a *App) sendUpdatedRoleEvent(role *model.Role) {
	message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_ROLE_UPDATED, "", "", "", nil)
	message.Add("role", role.ToJson())

	a.Go(func() {
		a.Publish(message)
	})
}
