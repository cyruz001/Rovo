package handler

import (
	"context"
	"strconv"

	"goServer/internal/dto"
	"goServer/internal/model"
	"goServer/internal/service"

	"github.com/gofiber/fiber/v3"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(ps *service.PostService) *PostHandler {
	return &PostHandler{postService: ps}
}

// CreatePost creates a new post
func (h *PostHandler) CreatePost(c fiber.Ctx) error {
	userID := c.Locals("sub").(string)
	var req dto.CreatePostReq

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	post, err := h.postService.CreatePost(context.Background(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(postToRes(post))
}

// GetPost retrieves a post by ID
func (h *PostHandler) GetPost(c fiber.Ctx) error {
	postID := c.Params("id")
	currentUserID := c.Locals("sub")

	post, err := h.postService.GetPostByID(context.Background(), postID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
	}

	res := postToRes(post)
	if currentUserID != nil {
		res.IsLiked, _ = h.postService.IsPostLiked(context.Background(), currentUserID.(string), postID)
		res.IsReposted, _ = h.postService.IsPostReposted(context.Background(), currentUserID.(string), postID)
	}

	return c.JSON(res)
}

// UpdatePost updates a post
func (h *PostHandler) UpdatePost(c fiber.Ctx) error {
	userID := c.Locals("sub").(string)
	postID := c.Params("id")
	var req dto.UpdatePostReq

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	post, err := h.postService.UpdatePost(context.Background(), postID, userID, req.Text)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(postToRes(post))
}

// DeletePost deletes a post
func (h *PostHandler) DeletePost(c fiber.Ctx) error {
	userID := c.Locals("sub").(string)
	postID := c.Params("id")

	if err := h.postService.DeletePost(context.Background(), postID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "post deleted"})
}

// GetFeed retrieves user's feed
func (h *PostHandler) GetFeed(c fiber.Ctx) error {
	userID := c.Locals("sub").(string)
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	posts, err := h.postService.GetFeed(context.Background(), userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.PostRes, len(posts))
	for i, p := range posts {
		r := postToRes(&p)
		r.IsLiked, _ = h.postService.IsPostLiked(context.Background(), userID, p.ID)
		r.IsReposted, _ = h.postService.IsPostReposted(context.Background(), userID, p.ID)
		res[i] = r
	}

	return c.JSON(res)
}

// GetUserTimeline retrieves user's timeline
func (h *PostHandler) GetUserTimeline(c fiber.Ctx) error {
	username := c.Params("username")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	currentUserID := c.Locals("sub")

	posts, err := h.postService.GetUserTimeline(context.Background(), username, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.PostRes, len(posts))
	for i, p := range posts {
		r := postToRes(&p)
		if currentUserID != nil {
			r.IsLiked, _ = h.postService.IsPostLiked(context.Background(), currentUserID.(string), p.ID)
			r.IsReposted, _ = h.postService.IsPostReposted(context.Background(), currentUserID.(string), p.ID)
		}
		res[i] = r
	}

	return c.JSON(res)
}

// LikePost likes a post
func (h *PostHandler) LikePost(c fiber.Ctx) error {
	userID := c.Locals("sub").(string)
	postID := c.Params("id")

	if err := h.postService.LikePost(context.Background(), userID, postID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "post liked"})
}

// UnlikePost unlikes a post
func (h *PostHandler) UnlikePost(c fiber.Ctx) error {
	userID := c.Locals("sub").(string)
	postID := c.Params("id")

	if err := h.postService.UnlikePost(context.Background(), userID, postID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "post unliked"})
}

// RepostPost reposts a post
func (h *PostHandler) RepostPost(c fiber.Ctx) error {
	userID := c.Locals("sub").(string)
	postID := c.Params("id")

	if err := h.postService.RepostPost(context.Background(), userID, postID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "post reposted"})
}

// UndoRepost undoes a repost
func (h *PostHandler) UndoRepost(c fiber.Ctx) error {
	userID := c.Locals("sub").(string)
	postID := c.Params("id")

	if err := h.postService.UndoRepost(context.Background(), userID, postID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "repost undone"})
}

// GetPostLikes gets users who liked a post
func (h *PostHandler) GetPostLikes(c fiber.Ctx) error {
	postID := c.Params("id")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	users, err := h.postService.GetPostLikes(context.Background(), postID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.UserRes, len(users))
	for i, u := range users {
		res[i] = userToRes(&u)
	}

	return c.JSON(res)
}

// GetPostReposts gets users who reposted a post
func (h *PostHandler) GetPostReposts(c fiber.Ctx) error {
	postID := c.Params("id")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	users, err := h.postService.GetPostReposts(context.Background(), postID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.UserRes, len(users))
	for i, u := range users {
		res[i] = userToRes(&u)
	}

	return c.JSON(res)
}

// GetReplies gets replies to a post
func (h *PostHandler) GetReplies(c fiber.Ctx) error {
	postID := c.Params("id")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	currentUserID := c.Locals("sub")

	posts, err := h.postService.GetReplies(context.Background(), postID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.PostRes, len(posts))
	for i, p := range posts {
		r := postToRes(&p)
		if currentUserID != nil {
			r.IsLiked, _ = h.postService.IsPostLiked(context.Background(), currentUserID.(string), p.ID)
			r.IsReposted, _ = h.postService.IsPostReposted(context.Background(), currentUserID.(string), p.ID)
		}
		res[i] = r
	}

	return c.JSON(res)
}

// SearchPosts searches for posts
func (h *PostHandler) SearchPosts(c fiber.Ctx) error {
	query := c.Query("query")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	currentUserID := c.Locals("sub")

	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query is required"})
	}

	posts, err := h.postService.SearchPosts(context.Background(), query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.PostRes, len(posts))
	for i, p := range posts {
		r := postToRes(&p)
		if currentUserID != nil {
			r.IsLiked, _ = h.postService.IsPostLiked(context.Background(), currentUserID.(string), p.ID)
			r.IsReposted, _ = h.postService.IsPostReposted(context.Background(), currentUserID.(string), p.ID)
		}
		res[i] = r
	}

	return c.JSON(res)
}

// GetAllPosts gets all posts (admin)
func (h *PostHandler) GetAllPosts(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	posts, err := h.postService.GetFeed(context.Background(), "", limit, offset) // Get all posts
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := make([]dto.PostRes, len(posts))
	for i, p := range posts {
		res[i] = postToRes(&p)
	}

	return c.JSON(res)
}

// AdminDeletePost deletes a post (admin)
func (h *PostHandler) AdminDeletePost(c fiber.Ctx) error {
	postID := c.Params("id")

	if err := h.postService.DeletePost(context.Background(), postID, ""); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "post deleted"})
}

// Helper function to convert Post model to PostRes DTO
func postToRes(p *model.Post) dto.PostRes {
	media := make([]dto.MediaRes, len(p.Media))
	for i, m := range p.Media {
		media[i] = dto.MediaRes{
			ID:        m.ID,
			URL:       m.URL,
			MediaType: m.MediaType,
			Position:  m.Position,
		}
	}

	return dto.PostRes{
		ID:           p.ID,
		UserID:       p.UserID,
		User:         userToRes(&p.User),
		Text:         p.Text,
		CharCount:    p.CharCount,
		ReplyTo:      p.ReplyTo,
		IsQuote:      p.IsQuote,
		QuotedPostID: p.QuotedTweetID,
		Media:        media,
		CreatedAt:    p.CreatedAt.String(),
		UpdatedAt:    p.UpdatedAt.String(),
	}
}
