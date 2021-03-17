package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"go.uber.org/zap"

	"github.com/spf13/pflag"

	chatmapper "github.com/arykalin/metaver-topos/chat_mapper"
	"github.com/arykalin/metaver-topos/mailer"
	"github.com/arykalin/metaver-topos/notifier"
	sheet "github.com/arykalin/metaver-topos/scheet"
	"github.com/arykalin/metaver-topos/users"
)

type Config struct {
	ChatMapSheetID   string `yaml:"chat_logic_sheet"`
	AnswersSheet1ID  string `yaml:"answers_sheet_1"`
	AnswersSheet2ID  string `yaml:"answers_sheet_2"`
	MailUser         string `yaml:"mail_user"`
	MailPassword     string `yaml:"mail_password"`
	MailHost         string `yaml:"mail_host"`
	MailPort         string `yaml:"mail_port"`
	MailDebugAddress string `yaml:"mail_debug_address"`
	MailCCAddress    string `yaml:"mail_cc_address"`
}

func main() {
	pathConfig := pflag.StringP("path", "c", "./config.yml", "path to config file")
	help := pflag.BoolP("help", "h", false, "show help")
	pflag.Parse()

	b, err := ioutil.ReadFile(*pathConfig)
	if err != nil {
		log.Fatalf("can't read file")
	}

	if *help {
		pflag.PrintDefaults()
		os.Exit(0)
	}

	var config Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		log.Fatalf("can't unmarshal config: %s", err)
	}

	sLoggerConfig := zap.NewDevelopmentConfig()
	sLoggerConfig.DisableStacktrace = true
	sLoggerConfig.DisableCaller = true
	sLogger, err := sLoggerConfig.Build()
	if err != nil {
		panic(err)
	}
	logger := sLogger.Sugar()

	s, err := sheet.NewSheetService(logger)
	if err != nil {
		log.Fatalf("failed to init sheet client: %s", err)
	}

	// form map for mailer
	chatMapSheetID := config.ChatMapSheetID
	chatMapSpreadsheet, err := s.GetSheet(chatMapSheetID)
	if err != nil {
		log.Fatalf("failed to get sheet data: %s", err)
	}
	chatMapSheet, err := chatMapSpreadsheet.SheetByIndex(0)
	if err != nil {
		log.Fatalf("failed to get sheet data: %s", err)
	}
	mapper := chatmapper.NewChatMapper(chatMapSheet)
	j, _ := json.MarshalIndent(mapper.GetMap(), "", "    ")
	fmt.Printf("got chat map:\n%s\n", j)

	// get answers
	formUsers := users.NewUsers()

	// Add users from first form
	spreadsheet1, err := s.GetSheet(config.AnswersSheet1ID)
	if err != nil {
		logger.Fatalf("failed to get sheet data: %s", err)
	}
	sheet1, err := spreadsheet1.SheetByIndex(0)
	if err != nil {
		logger.Fatalf("failed to get sheet data: %s", err)
	}
	sheet1Config := users.SheetConfig{
		TrackIdx:    1,
		MailIdx:     2,
		HaveTeam:    true,
		UserTypeIdx: nil,
	}
	err = formUsers.AddUsers(sheet1, &sheet1Config)
	if err != nil {
		logger.Fatalf("failed to make users map: %s", err)
	}
	logger.Debugw("users after first form", "total", len(formUsers.GetUsers()))

	// Add users from second form
	spreadsheet2, err := s.GetSheet(config.AnswersSheet2ID)
	if err != nil {
		logger.Fatalf("failed to get sheet data: %s", err)
	}
	sheet2, err := spreadsheet2.SheetByIndex(0)
	if err != nil {
		logger.Fatalf("failed to get sheet data: %s", err)
	}

	mentorIdx := 8
	sheet2Config := users.SheetConfig{
		TrackIdx:    1,
		MailIdx:     6,
		HaveTeam:    false,
		UserTypeIdx: &mentorIdx,
	}
	err = formUsers.AddUsers(sheet2, &sheet2Config)
	if err != nil {
		logger.Fatalf("failed to make users map: %s", err)
	}
	logger.Debugw("users after second form", "total", len(formUsers.GetUsers()))
	err = formUsers.DumpUsers()
	if err != nil {
		logger.Fatalf("failed to dump users map: %s", err)
	}

	newMailer := mailer.NewMailer(
		logger,
		config.MailUser,
		config.MailPassword,
		config.MailHost,
		config.MailDebugAddress,
		config.MailCCAddress,
	)
	n := notifier.NewNotifier(logger, newMailer)
	err = n.Notify(mapper.GetMap(), formUsers.GetUsers())
	if err != nil {
		logger.Fatalf("error notify: %s", err)
	}
}
