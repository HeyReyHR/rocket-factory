package converter

import (
	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
	repoModel "github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
)

func ConvertUserInfoServiceToRepo(info model.UserInfo) repoModel.UserInfo {
	return repoModel.UserInfo{
		Email:               info.Email,
		Login:               info.Login,
		NotificationMethods: ConvertNotificationMethodsServiceToRepo(info.NotificationMethods),
	}
}

func ConvertNotificationMethodsServiceToRepo(methods []model.NotificationMethod) []repoModel.NotificationMethod {
	var repoMethods []repoModel.NotificationMethod
	for _, method := range methods {
		repoMethods = append(repoMethods, repoModel.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}
	
	return repoMethods
}
