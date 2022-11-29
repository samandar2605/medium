package v1

// import (
// 	"net/http"

// 	"net/smtp"

// 	"github.com/gin-gonic/gin"
// 	"github.com/post/api/models"
// )

// // @Router /emails{email} [post]
// // @Summary Create a emails
// // @Description Create a emails
// // @Tags email
// // @Accept json
// // @Produce json
// // @Param email path string true "from_email"
// // @Param api_key path string true "your_api_key"
// // @Param email body models.CreateEmail true "email"
// // @Success 201 {object} models.Email
// // @Failure 400 {object} models.ErrorResponse
// // @Failure 500 {object} models.ErrorResponse
// func (h *handlerV1) SendEmail(c *gin.Context) {
// 	var (
// 		req models.CreateEmail
// 	)

// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
// 			Error: err.Error(),
// 		})
// 		return
// 	}

// 	from := req.YourEmail
// 	password := "jrqgvejmkblgjxft"

// 	toEmailAddress := "t.mannonov@gmail.com"
// 	to := []string{toEmailAddress}

// 	host := "smtp.gmail.com"
// 	port := "587"
// 	address := host + ":" + port

// 	subject := "Samandar\n"
// 	body := "Ustoz bugun nima dars o'tamiz"
// 	message := []byte(subject + body)

// 	auth := smtp.PlainAuth("", from, password, host)

// 	err = smtp.SendMail(address, auth, from, to, message)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
// 			Error: err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, models.Email{})
// }
