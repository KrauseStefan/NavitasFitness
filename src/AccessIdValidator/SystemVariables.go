package AccessIdValidator

import (
	"golang.org/x/net/context"

	"SystemSettingDAO"
	"google.golang.org/appengine/log"
)

const (
	fitnessAccessIdsPathSettingKey = "fitnessAccessIdsPath"
	defaultFitnessAccessIdsPath    = "/AccessIds/AccessIds.csv"

	fitnessAccessListPathSettingKey = "fitnessAccessListPath"
	defaultFitnessAccessListPath    = "/FitnessAccessList/FitnessAccessList.csv"
)

func GetAccessIdPath(ctx context.Context) string {
	_, value, err := SystemSettingDAO.GetSetting(ctx, fitnessAccessIdsPathSettingKey)
	if err != nil {
		log.Errorf(ctx, err.Error())
	}

	if value == "" {
		value = defaultFitnessAccessIdsPath
		if err := SystemSettingDAO.PersistSetting(ctx, fitnessAccessIdsPathSettingKey, value); err != nil {
			log.Errorf(ctx, err.Error())
		}
	}

	return value
}

func GetAccessListPath(ctx context.Context) string {
	_, value, err := SystemSettingDAO.GetSetting(ctx, fitnessAccessListPathSettingKey)
	if err != nil {
		log.Infof(ctx, err.Error())
	}

	if value == "" {
		value = defaultFitnessAccessListPath
		if err := SystemSettingDAO.PersistSetting(ctx, fitnessAccessListPathSettingKey, value); err != nil {
			log.Infof(ctx, err.Error())
		}
	}

	return value
}
