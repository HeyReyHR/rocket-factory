package converter

import (
	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
	repoModel "github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
	commonV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/common/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func ConvertUserServiceToApi(user model.User) *commonV1.User {
	return &commonV1.User{
		Uuid:      user.Uuid,
		Info:      ConvertUserInfoServiceToApi(user.UserInfo),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ConvertUserInfoServiceToApi(info model.UserInfo) *commonV1.UserInfo {
	return &commonV1.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		NotificationMethods: ConvertNotificationMethodsServiceToApi(info.NotificationMethods),
	}
}

func ConvertNotificationMethodsServiceToApi(methods []model.NotificationMethod) []*commonV1.NotificationMethod {
	var apiMethods []*commonV1.NotificationMethod
	for _, method := range methods {
		apiMethods = append(apiMethods, &commonV1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return apiMethods
}

func ConvertSessionServiceToApi(session model.Session) *commonV1.Session {
	return &commonV1.Session{
		Uuid:      session.Uuid,
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: timestamppb.New(session.UpdatedAt),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}

func ConvertUserInfoApiToService(info *commonV1.UserInfo) model.UserInfo {
	return model.UserInfo{
		Login:               info.GetLogin(),
		Email:               info.GetEmail(),
		NotificationMethods: ConvertNotificationMethodsApiToService(info.GetNotificationMethods()),
	}
}

func ConvertNotificationMethodsApiToService(methods []*commonV1.NotificationMethod) []model.NotificationMethod {
	var serviceMethods []model.NotificationMethod
	for _, method := range methods {
		serviceMethods = append(serviceMethods, model.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return serviceMethods
}
