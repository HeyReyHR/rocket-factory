package converter

import (
	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
	repoModel "github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
)

func ConvertUserRepoToService(user repoModel.User) model.User {
	return model.User{
		Uuid:      user.Uuid,
		UserInfo:  ConvertUserInfoRepoToService(user.UserInfo),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ConvertUserInfoRepoToService(info repoModel.UserInfo) model.UserInfo {
	return model.UserInfo{
		Email:               info.Email,
		Login:               info.Login,
		NotificationMethods: ConvertNotificationMethodsRepoToService(info.NotificationMethods),
	}
}

func ConvertNotificationMethodsRepoToService(methods []repoModel.NotificationMethod) []model.NotificationMethod {
	var repoMethods []model.NotificationMethod
	for _, method := range methods {
		repoMethods = append(repoMethods, model.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return repoMethods
}

func ConvertSessionRepoToService(session repoModel.Session) model.Session {
	return model.Session{
		Uuid:      session.Uuid,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
		ExpiresAt: session.ExpiresAt,
	}
}
