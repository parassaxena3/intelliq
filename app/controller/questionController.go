package controller

import (
	"intelliq/app/common"
	utility "intelliq/app/common"
	"intelliq/app/dto"
	"intelliq/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

//FindQuestion get question based on quesId
func FindQuestion(ctx *gin.Context) {
	groupCode := ctx.Param("groupCode")
	quesID := ctx.Param("quesId")
	res := service.FetchOneQuestion(groupCode, quesID)
	ctx.JSON(http.StatusOK, res)
}

//GetQuestionsFromBank get all approved questions from bank
func GetQuestionsFromBank(ctx *gin.Context) {
	var requestDto dto.QuesRequestDto
	err := ctx.BindJSON(&requestDto)
	if err != nil {
		res := utility.GetErrorResponse(common.MSG_BAD_INPUT)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := service.FetchApprovedQuestions(&requestDto)
	ctx.JSON(http.StatusOK, res)
}

//GetQuestionSuggestions returns question suggestions on typing new ques
func GetQuestionSuggestions(ctx *gin.Context) {
	var quesCriteriaDto dto.QuestionCriteriaDto
	err := ctx.BindJSON(&quesCriteriaDto)
	if err != nil {
		res := utility.GetErrorResponse(common.MSG_BAD_INPUT)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := service.FetchQuestionSuggestions(&quesCriteriaDto)
	ctx.JSON(http.StatusOK, res)
}

//GetFilteredQuestions filters questions as per criteria
func GetFilteredQuestions(ctx *gin.Context) {
	var quesCriteriaDto dto.QuestionCriteriaDto
	err := ctx.BindJSON(&quesCriteriaDto)
	if err != nil {
		res := utility.GetErrorResponse(common.MSG_BAD_INPUT)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := service.FilterQuestions(&quesCriteriaDto)
	ctx.JSON(http.StatusOK, res)
}

//RemoveObsoleteQuestions removes all obsolete questions
func RemoveObsoleteQuestions(ctx *gin.Context) {
	groupCode := ctx.Param("groupCode")
	res := service.RemoveObsoleteQuestions(groupCode)
	ctx.JSON(http.StatusOK, res)
}
