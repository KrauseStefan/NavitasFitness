package AccessIdValidator

import (
	"golang.org/x/net/context"

	"github.com/KrauseStefan/NavitasFitness/SystemSettingDAO"
	log "github.com/KrauseStefan/NavitasFitness/logger"
)

const (
	// Below prefix is the default and is used for testing,
	// in the production environment this is overwriten to be the real path
	defaultTestDropboxPrefix       = "/Test/3210 Navitas ADK - Fitness"
	fitnessAccessIdsPathSettingKey = "fitnessAccessIdsPath"
	defaultFitnessAccessIdsPath    = defaultTestDropboxPrefix + "/AccessIds/AccessIds.csv"

	fitnessAccessListPathSettingKey = "fitnessAccessListPath"
	defaultFitnessAccessListPath    = defaultTestDropboxPrefix + "/FitnessAccessList/FitnessAccessList.csv"

	paypallValidationEmailSettingKey = "paypallValidationEmail"
	defaultPaypallValidationEmail    = "navitasShop2@mail.dk:gpmac_1231902686_biz@paypal.com:navitasShop@mail.dk"
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

func GetPaypalValidationEmail(ctx context.Context) string {
	_, value, err := SystemSettingDAO.GetSetting(ctx, paypallValidationEmailSettingKey)
	if err != nil {
		log.Infof(ctx, err.Error())
	}

	if value == "" {
		value = defaultPaypallValidationEmail
		if err := SystemSettingDAO.PersistSetting(ctx, paypallValidationEmailSettingKey, value); err != nil {
			log.Infof(ctx, err.Error())
		}
	}

	return value
}
