package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	dao "github.com/ninjadotorg/handshake-telegram/dao"
	models "github.com/ninjadotorg/handshake-telegram/models"
	"github.com/urfave/cli"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	app           *cli.App
	chatMemberDAO = dao.ChatMemberDao{}
)

func init() {
	// Initialise a CLI app
	app = cli.NewApp()
	app.Name = "ninja ethereum"
	app.Usage = "ninja ethereum"
	app.Author = "hieuqautonomous"
	app.Email = "hieu.q@autonomous.nyc"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c",
			Value: "",
			Usage: "Path to a configuration file",
		},
	}
}

func main() {
	app.Commands = []cli.Command{
		{
			Name:  "worker",
			Usage: "launch worker",
			Action: func(c *cli.Context) error {
				return workerApp()
			},
		},
		{
			Name:  "service",
			Usage: "launch service",
			Action: func(c *cli.Context) error {
				return serviceApp()
			},
		},
	}
	// Run the CLI app
	if err := app.Run(os.Args); err != nil {
		log.Println("error", err)
	}
	select {}
}

func workerApp() error {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		panic(errors.New("env is invalid"))
	}
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.NewChatMember != nil {
			newChatMember := update.Message.NewChatMember
			chatMember := chatMemberDAO.GetByFilter(update.Message.Chat.ID, int64(newChatMember.ID))
			if chatMember.ID <= 0 {
				chatMember = models.ChatMember{}
				chatMember.ChatID = update.Message.Chat.ID
				chatMember.UserID = int64(newChatMember.ID)
				chatMember.UserName = newChatMember.UserName
				chatMember.FirstName = newChatMember.FirstName
				chatMember.LastName = newChatMember.LastName
				chatMember, err = chatMemberDAO.Create(chatMember, nil)
				if err != nil {
					log.Println(err)
				}
			}
		}
		if update.Message.LeftChatMember != nil {
			leftChatMember := update.Message.LeftChatMember
			chatMember := chatMemberDAO.GetByFilter(update.Message.Chat.ID, int64(leftChatMember.ID))
			if chatMember.ID <= 0 {
				chatMember, err = chatMemberDAO.Delete(chatMember, nil)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	return nil
}

func serviceApp() error {
	// Logger
	logFile, err := os.OpenFile("logs/autonomous_service.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(gin.DefaultWriter) // You may need this
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	router := gin.Default()
	router.Use(Logger())
	router.Use(AuthorizeMiddleware())
	index := router.Group("/")
	{
		index.GET("/", func(c *gin.Context) {
			result := map[string]interface{}{
				"status":  1,
				"message": "Telegram Service API",
			}
			c.JSON(http.StatusOK, result)
		})
		index.GET("/chat_member", func(c *gin.Context) {
			chatID, _ := strconv.ParseInt(c.Query("chat_id"), 10, 64)
			userName := c.Query("user_name")
			chatMember := chatMemberDAO.GetByUserName(chatID, userName)
			if chatMember.ID > 0 {
				result := map[string]interface{}{
					"status":  1,
					"message": "OK",
					"data":    chatMember,
				}
				c.JSON(http.StatusOK, result)
			} else {
				result := map[string]interface{}{
					"status":  1,
					"message": "OK",
				}
				c.JSON(http.StatusOK, result)
			}
		})
	}

	router.Run(":8080")

	return nil
}

func Logger() gin.HandlerFunc {
	return func(context *gin.Context) {
		t := time.Now()
		context.Next()
		status := context.Writer.Status()
		latency := time.Since(t)
		log.Print("Request: " + context.Request.URL.String() + " | " + context.Request.Method + " - Status: " + strconv.Itoa(status) + " - " +
			latency.String())
	}
}

func AuthorizeMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		userID, _ := strconv.ParseInt(context.GetHeader("Uid"), 10, 64)
		if userID <= 0 {
			context.JSON(http.StatusOK, gin.H{"status": 0, "message": "User is not authorized"})
			context.Abort()
			return
		}
		context.Set("UserID", userID)
		context.Next()
	}
}
