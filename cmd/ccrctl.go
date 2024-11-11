package cmd

import (
	"ccrctl/pkg/logger"
	"ccrctl/pkg/migrate"
	"os"
)

var (
	Version   string
	BuildTime string
)

func ccrctl() {
	st, err := os.Lstat(os.Args[0])
	if err != nil {
		logger.Logger.Errorf("os.Lastat error: %v", err)
	}

	logger.Logger.Infof("===========================================================")
	logger.Logger.Infof("|  Server Name : %-40s |", st.Name())
	logger.Logger.Infof("|  Build Time  : %-40s |", BuildTime)
	logger.Logger.Infof("|  Version     : %-39s |", Version)
	logger.Logger.Infof("===========================================================")
	migrate.Run()
}
