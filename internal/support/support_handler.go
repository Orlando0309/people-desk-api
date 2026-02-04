package support
import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Create(c *gin.Context) {
	var input CreateSupportInput
	if err:= c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var support= &Support{
		Message: input.Message,
		Email: input.Email,
	}
	created,err:= h.repo.Create(c.Request.Context(), support)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}