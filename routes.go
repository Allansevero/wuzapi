package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type Middleware = alice.Constructor

func (s *server) routes() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	var routerLog zerolog.Logger
	if *logType == "json" {
		routerLog = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Str("role", filepath.Base(os.Args[0])).
			Str("host", *address).
			Logger()
	} else {
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    !*colorOutput,
		}
		routerLog = zerolog.New(output).
			With().
			Timestamp().
			Str("role", filepath.Base(os.Args[0])).
			Str("host", *address).
			Logger()
	}

	s.router.Handle("/health", s.GetHealth()).Methods("GET")
	s.router.Handle("/healthz", s.healthCheckHandler()).Methods("GET")

	// Authentication routes (public)
	s.router.Handle("/auth/register", s.Register()).Methods("POST")
	s.router.Handle("/auth/login", s.Login()).Methods("POST")
	s.router.Handle("/auth/logout", s.SystemLogout()).Methods("POST")

	// Instance management routes (authenticated system users)
	userRoutes := s.router.PathPrefix("/my").Subrouter()
	userRoutes.Use(s.authSystemUser)
	
	// Profile routes
	userRoutes.Handle("/profile", s.GetMyProfile()).Methods("GET")
	userRoutes.Handle("/profile", s.UpdateMyProfile()).Methods("PUT")
	
	// Instance routes
	userRoutes.Handle("/instances", s.ListMyInstances()).Methods("GET")
	userRoutes.Handle("/instances", s.CreateMyInstance()).Methods("POST")
	userRoutes.Handle("/instances/{id}", s.GetMyInstance()).Methods("GET")
	userRoutes.Handle("/instances/{id}", s.UpdateMyInstance()).Methods("PUT")
	userRoutes.Handle("/instances/{id}", s.DeleteMyInstance()).Methods("DELETE")
	
	// Subscription routes (authenticated system users)
	userRoutes.Handle("/subscription", s.GetUserSubscriptionHandler()).Methods("GET")
	userRoutes.Handle("/subscription", s.UpdateUserSubscriptionHandler()).Methods("PUT")
	userRoutes.Handle("/plans", s.GetPlansHandler()).Methods("GET")

	adminRoutes := s.router.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(s.authadmin)
	adminRoutes.Handle("/users", s.ListUsers()).Methods("GET")
	adminRoutes.Handle("/users/{id}", s.ListUsers()).Methods("GET")
	adminRoutes.Handle("/users", s.AddUser()).Methods("POST")
	adminRoutes.Handle("/users/{id}", s.EditUser()).Methods("PUT")
	adminRoutes.Handle("/users/{id}", s.DeleteUser()).Methods("DELETE")
	adminRoutes.Handle("/users/{id}/full", s.DeleteUserComplete()).Methods("DELETE")
	
	// Admin route for getting all instances with names and destination numbers
	adminRoutes.Handle("/instances", s.ListInstancesForAdmin()).Methods("GET")
	
	// Admin route for pushing chat history to webhook
	adminRoutes.Handle("/chat/history/push", s.PushHistoryToWebhook()).Methods("POST")

	c := alice.New()
	c = c.Append(s.authalice)
	c = c.Append(hlog.NewHandler(routerLog))

	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Str("userid", r.Context().Value("userinfo").(Values).Get("Id")).
			Msg("Got API Request")
	}))

	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	// Chain with subscription check for chat endpoints
	cWithSub := c.Append(s.checkSubscriptionMiddleware)

	s.router.Handle("/session/connect", c.Then(s.Connect())).Methods("POST")
	s.router.Handle("/session/disconnect", c.Then(s.Disconnect())).Methods("POST")
	s.router.Handle("/session/logout", c.Then(s.Logout())).Methods("POST")
	s.router.Handle("/session/status", c.Then(s.GetStatus())).Methods("GET")
	s.router.Handle("/session/qr", c.Then(s.GetQR())).Methods("GET")
	s.router.Handle("/session/pairphone", c.Then(s.PairPhone())).Methods("POST")
	s.router.Handle("/session/history", c.Then(s.RequestHistorySync())).Methods("GET")

	s.router.Handle("/webhook", c.Then(s.SetWebhook())).Methods("POST")
	s.router.Handle("/webhook", c.Then(s.GetWebhook())).Methods("GET")
	s.router.Handle("/webhook", c.Then(s.DeleteWebhook())).Methods("DELETE")
	s.router.Handle("/webhook", c.Then(s.UpdateWebhook())).Methods("PUT")

	s.router.Handle("/session/proxy", c.Then(s.SetProxy())).Methods("POST")
	s.router.Handle("/session/history", c.Then(s.SetHistory())).Methods("POST")
	
	// Destination number configuration
	s.router.Handle("/session/destination-number", c.Then(s.SetDestinationNumber())).Methods("POST")
	s.router.Handle("/session/destination-number", c.Then(s.GetDestinationNumber())).Methods("GET")
	


	s.router.Handle("/session/s3/config", c.Then(s.ConfigureS3())).Methods("POST")
	s.router.Handle("/session/s3/config", c.Then(s.GetS3Config())).Methods("GET")
	s.router.Handle("/session/s3/config", c.Then(s.DeleteS3Config())).Methods("DELETE")
	s.router.Handle("/session/s3/test", c.Then(s.TestS3Connection())).Methods("POST")

	s.router.Handle("/session/hmac/config", c.Then(s.ConfigureHmac())).Methods("POST")
	s.router.Handle("/session/hmac/config", c.Then(s.GetHmacConfig())).Methods("GET")
	s.router.Handle("/session/hmac/config", c.Then(s.DeleteHmacConfig())).Methods("DELETE")

	// Chat endpoints with subscription check
	s.router.Handle("/chat/send/text", cWithSub.Then(s.SendMessage())).Methods("POST")
	s.router.Handle("/chat/delete", cWithSub.Then(s.DeleteMessage())).Methods("POST")
	s.router.Handle("/chat/send/image", cWithSub.Then(s.SendImage())).Methods("POST")
	s.router.Handle("/chat/send/audio", cWithSub.Then(s.SendAudio())).Methods("POST")
	s.router.Handle("/chat/send/document", cWithSub.Then(s.SendDocument())).Methods("POST")
	//	s.router.Handle("/chat/send/template", cWithSub.Then(s.SendTemplate())).Methods("POST")
	s.router.Handle("/chat/send/video", cWithSub.Then(s.SendVideo())).Methods("POST")
	s.router.Handle("/chat/send/sticker", cWithSub.Then(s.SendSticker())).Methods("POST")
	s.router.Handle("/chat/send/location", cWithSub.Then(s.SendLocation())).Methods("POST")
	s.router.Handle("/chat/send/contact", cWithSub.Then(s.SendContact())).Methods("POST")
	s.router.Handle("/chat/react", cWithSub.Then(s.React())).Methods("POST")
	s.router.Handle("/chat/send/buttons", cWithSub.Then(s.SendButtons())).Methods("POST")
	s.router.Handle("/chat/send/list", cWithSub.Then(s.SendList())).Methods("POST")
	s.router.Handle("/chat/send/poll", cWithSub.Then(s.SendPoll())).Methods("POST")
	s.router.Handle("/chat/send/edit", cWithSub.Then(s.SendEditMessage())).Methods("POST")
	s.router.Handle("/chat/history", cWithSub.Then(s.GetHistory())).Methods("GET")
	s.router.Handle("/chat/history/push", cWithSub.Then(s.PushHistoryToWebhookByInstance())).Methods("POST")

	s.router.Handle("/status/set/text", c.Then(s.SetStatusMessage())).Methods("POST")

	s.router.Handle("/call/reject", c.Then(s.RejectCall())).Methods("POST")

	s.router.Handle("/user/presence", c.Then(s.SendPresence())).Methods("POST")
	s.router.Handle("/user/info", c.Then(s.GetUser())).Methods("POST")
	s.router.Handle("/user/check", c.Then(s.CheckUser())).Methods("POST")
	s.router.Handle("/user/avatar", c.Then(s.GetAvatar())).Methods("POST")
	s.router.Handle("/user/contacts", c.Then(s.GetContacts())).Methods("GET")
	s.router.Handle("/user/lid/{jid}", c.Then(s.GetUserLID())).Methods("GET")

	s.router.Handle("/chat/presence", c.Then(s.ChatPresence())).Methods("POST")
	s.router.Handle("/chat/markread", c.Then(s.MarkRead())).Methods("POST")
	s.router.Handle("/chat/downloadimage", c.Then(s.DownloadImage())).Methods("POST")
	s.router.Handle("/chat/downloadvideo", c.Then(s.DownloadVideo())).Methods("POST")
	s.router.Handle("/chat/downloadaudio", c.Then(s.DownloadAudio())).Methods("POST")
	s.router.Handle("/chat/downloaddocument", c.Then(s.DownloadDocument())).Methods("POST")

	s.router.Handle("/group/create", c.Then(s.CreateGroup())).Methods("POST")
	s.router.Handle("/group/list", c.Then(s.ListGroups())).Methods("GET")
	s.router.Handle("/group/info", c.Then(s.GetGroupInfo())).Methods("GET")
	s.router.Handle("/group/invitelink", c.Then(s.GetGroupInviteLink())).Methods("GET")
	s.router.Handle("/group/photo", c.Then(s.SetGroupPhoto())).Methods("POST")
	s.router.Handle("/group/photo/remove", c.Then(s.RemoveGroupPhoto())).Methods("POST")
	s.router.Handle("/group/leave", c.Then(s.GroupLeave())).Methods("POST")
	s.router.Handle("/group/name", c.Then(s.SetGroupName())).Methods("POST")
	s.router.Handle("/group/topic", c.Then(s.SetGroupTopic())).Methods("POST")
	s.router.Handle("/group/announce", c.Then(s.SetGroupAnnounce())).Methods("POST")
	s.router.Handle("/group/locked", c.Then(s.SetGroupLocked())).Methods("POST")
	s.router.Handle("/group/ephemeral", c.Then(s.SetDisappearingTimer())).Methods("POST")
	s.router.Handle("/group/join", c.Then(s.GroupJoin())).Methods("POST")
	s.router.Handle("/group/inviteinfo", c.Then(s.GetGroupInviteInfo())).Methods("POST")
	s.router.Handle("/group/updateparticipants", c.Then(s.UpdateGroupParticipants())).Methods("POST")

	s.router.Handle("/newsletter/list", c.Then(s.ListNewsletter())).Methods("GET")

	// Rotas específicas para arquivos estáticos antes do PathPrefix("/")
	s.router.Handle("/user-login.html", http.FileServer(http.Dir(exPath + "/static/")))
	s.router.PathPrefix("/login/").Handler(http.StripPrefix("/login/", http.FileServer(http.Dir(exPath + "/static/login/"))))
	s.router.PathPrefix("/dashboard/").Handler(http.StripPrefix("/dashboard/", http.FileServer(http.Dir(exPath + "/static/dashboard/"))))
	s.router.PathPrefix("/api/").Handler(http.StripPrefix("/api/", http.FileServer(http.Dir(exPath + "/static/api/"))))

	// Rota genérica para outros arquivos estáticos
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir(exPath + "/static/")))
}
