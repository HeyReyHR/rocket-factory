package converter

import (
	"strings"

	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
	repoModel "github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
)

func ConvertUserRepoToService(user repoModel.User) model.User {
	return model.User{
		Uuid:           user.Uuid,
		AdditionalInfo: ConvertUserInfoRepoToService(user.AdditionalInfo),
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func ConvertUserInfoRepoToService(info repoModel.AdditionalInfo) model.AdditionalInfo {
	return model.AdditionalInfo{
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

func ConvertUserDbToRepo(dto *repoModel.UserDbDto) (repoModel.User, error) {
	methods := make([]repoModel.NotificationMethod, 0, len(dto.NotificationMethods))
	for _, nm := range dto.NotificationMethods {
		parts := strings.SplitN(nm, ":", 2)
		if len(parts) == 2 {
			methods = append(methods, repoModel.NotificationMethod{
				ProviderName: parts[0],
				Target:       parts[1],
			})
		}
	}

	user := repoModel.User{
		Uuid:      dto.Uuid,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		AdditionalInfo: repoModel.AdditionalInfo{
			Login:               dto.Login,
			Email:               dto.Email,
			PasswordHash:        dto.PasswordHash,
			NotificationMethods: methods,
		},
	}

	return user, nil
}
