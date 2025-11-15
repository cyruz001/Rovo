package router

import (
	"time"

	"goServer/internal/config"
	"goServer/internal/handler"
	"goServer/internal/middleware"
	"goServer/internal/repository"
	"goServer/internal/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, cfg config.Config) {
	// Dependency Injection - Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	rateLimitRepo := repository.NewRateLimitRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// Dependency Injection - Services
	userSvc := service.NewUserService(*userRepo)
	postSvc := service.NewPostService(*postRepo, *userRepo)
	rateLimitSvc := service.NewRateLimitService(*rateLimitRepo)
	notificationSvc := service.NewNotificationService(*notificationRepo, *userRepo)

	// Dependency Injection - Handlers
	authHandler := handler.NewAuthHandler(userSvc, cfg)
	userHandler := handler.NewUserHandler(userSvc, notificationSvc)
	postHandler := handler.NewPostHandler(postSvc)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// ============ PUBLIC ROUTES ============
	// Authentication
	v1.Post("/auth/register", authHandler.Register)
	v1.Post("/auth/login", authHandler.Login)

	// Public Search (MUST BE BEFORE :username route)
	v1.Get("/users/search", userHandler.SearchUsers)

	// Public User Info (SPECIFIC ROUTES BEFORE WILDCARD)
	v1.Get("/users/:username/followers", userHandler.GetFollowers)
	v1.Get("/users/:username/following", userHandler.GetFollowing)
	v1.Get("/users/:username", userHandler.GetUserByUsername)

	// Public Posts
	v1.Get("/posts/:id/likes", postHandler.GetPostLikes)
	v1.Get("/posts/:id/reposts", postHandler.GetPostReposts)
	v1.Get("/posts/:id/replies", postHandler.GetReplies)
	v1.Get("/posts/:id", postHandler.GetPost)

	// ============ PROTECTED ROUTES (Requires JWT) ============
	protected := v1.Group("/", middleware.JWT(cfg.JWTSecret))

	// User Profile Management
	protected.Get("/users/me", userHandler.GetProfile)
	protected.Put("/users/me", userHandler.UpdateProfile)
	protected.Delete("/users/me", userHandler.DeleteAccount)

	protected.Get("/users/me/followers", userHandler.GetMyFollowers)
	protected.Get("/users/me/following", userHandler.GetMyFollowing)

	// User Relationships
	protected.Post("/users/:id/follow", userHandler.FollowUser)
	protected.Post("/users/:id/unfollow", userHandler.UnfollowUser)

	// Posts - Create & Manage
	protected.Post("/posts",
		middleware.RateLimit(rateLimitSvc, "create_post", 10, 15*time.Minute),
		postHandler.CreatePost)

	protected.Get("/posts/feed", postHandler.GetFeed)
	protected.Get("/posts/timeline/:username", postHandler.GetUserTimeline)

	// Generic post routes
	protected.Put("/posts/:id", postHandler.UpdatePost)
	protected.Delete("/posts/:id", postHandler.DeletePost)

	// Posts - Interactions
	protected.Post("/posts/:id/like",
		middleware.RateLimit(rateLimitSvc, "like_post", 50, 1*time.Minute),
		postHandler.LikePost)
	protected.Delete("/posts/:id/unlike", postHandler.UnlikePost)

	protected.Post("/posts/:id/repost",
		middleware.RateLimit(rateLimitSvc, "repost_post", 30, 1*time.Minute),
		postHandler.RepostPost)
	protected.Delete("/posts/:id/unrepost", postHandler.UndoRepost)

	// Notifications
	protected.Get("/notifications",
		middleware.RateLimit(rateLimitSvc, "get_notifications", 100, 1*time.Minute),
		userHandler.GetNotifications)
	protected.Put("/notifications/:id/read", userHandler.MarkNotificationAsRead)
	protected.Delete("/notifications/:id", userHandler.DeleteNotification)

	// Search
	protected.Get("/search/posts", postHandler.SearchPosts)
	protected.Get("/search/users", userHandler.SearchUsers)

	// ============ ADMIN ROUTES (Requires Admin Role) ============
	admin := protected.Group("/admin", middleware.RequireRole("ADMIN"))

	admin.Get("/users", userHandler.GetAllUsers)
	admin.Delete("/users/:id", userHandler.AdminDeleteUser)
	admin.Put("/users/:id/role", userHandler.UpdateUserRole)

	admin.Get("/posts", postHandler.GetAllPosts)
	admin.Delete("/posts/:id", postHandler.AdminDeletePost)

	admin.Get("/stats", userHandler.GetSystemStats)
}
