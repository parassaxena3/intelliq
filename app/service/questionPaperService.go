package service

import (
	"fmt"
	"intelliq/app/common"
	utility "intelliq/app/common"
	"intelliq/app/dto"
	"intelliq/app/enums"
	"intelliq/app/helper"
	"intelliq/app/model"
	"intelliq/app/repo"
	"sync"
	"time"

	"github.com/globalsign/mgo/bson"
)

//GenerateQuestionPaper generated question paper as per criteria provided
func GenerateQuestionPaper(criteriaDto *dto.QuestionCriteriaDto) *dto.AppResponseDto {
	errResponse := validateRequest(criteriaDto.GroupCode,
		criteriaDto.Subject, criteriaDto.Standard)
	if errResponse != nil {
		return errResponse
	}
	quesRepo := repo.NewQuestionRepository(criteriaDto.GroupCode)
	if quesRepo == nil {
		return utility.GetErrorResponse(common.MSG_UNATHORIZED_ACCESS)
	}
	criteriaDto.GenerateNatives()
	dbstart := time.Now()
	sectionQuesMap, err := quesRepo.FilterQuestionsForPaper(criteriaDto)
	if err != nil {
		fmt.Println(err.Error())
		errorMsg := utility.GetErrorMsg(err)
		if len(errorMsg) > 0 {
			return utility.GetErrorResponse(errorMsg)
		}
		return utility.GetErrorResponse(common.MSG_REQUEST_FAILED)
	}
	if len(sectionQuesMap) == 0 {
		return utility.GetErrorResponse(common.MSG_NO_RECORD)
	}
	fmt.Println("DB QUERY TIME := ", time.Since(dbstart))
	start := time.Now()
	helper.PrioritiseDifficultyList(criteriaDto.Difficulty)
	sectionCountMap := helper.GetSectionCountMap(criteriaDto.Length)
	sectionChannel := make(chan *dto.Section)
	for section, lvlMap := range sectionQuesMap { // go routine
		go helper.GetResultSectionQuesList(lvlMap, criteriaDto.Difficulty,
			section, sectionCountMap[section], criteriaDto.Sets, sectionChannel)
	}
	var quesSectionList []dto.Section
	for i := 0; i < len(sectionQuesMap); i++ {
		quesSectionList = append(quesSectionList, *<-sectionChannel)
	}
	close(sectionChannel)
	paperChannel := make(chan *dto.QuestionPaperDto)
	for currSet := 0; currSet < criteriaDto.Sets; currSet++ { // go routine
		go helper.GenerateQuestionPaper(quesSectionList,
			currSet, criteriaDto.Difficulty, paperChannel)
	}
	var questionPapers []dto.QuestionPaperDto
	for currSet := 0; currSet < criteriaDto.Sets; currSet++ {
		questionPapers = append(questionPapers, *<-paperChannel)
	}
	close(paperChannel)
	fmt.Println("TOTAL ALGO TIME := ", time.Since(start))
	return utility.GetSuccessResponse(questionPapers)
}

//SaveTestDetails saves template and question papers
func SaveTestDetails(testDto *dto.TestDto, saveAsDraft bool) *dto.AppResponseDto {
	var wg sync.WaitGroup
	tempChannel, paperChannel := make(chan string, 1), make(chan string, 1)
	go saveTemplate(testDto.Template, &wg, tempChannel)
	go saveTestPaper(testDto.TestPaper, saveAsDraft, &wg, paperChannel)
	wg.Wait()
	res := <-tempChannel + "\n" + <-paperChannel
	return utility.GetSuccessResponse(res)
}

func saveTemplate(template *model.Template, wg *sync.WaitGroup,
	tempChannel chan<- string) {
	wg.Add(1)
	defer cleanPanic(tempChannel, wg)
	if template == nil {
		tempChannel <- ""
	}
	template.Criteria512Hash = utility.GenerateHash(template.Criteria)
	if len(template.Criteria512Hash) == 0 {
		tempChannel <- common.MSG_CORRUPT_DATA
	}
	template.LastModifiedDate = time.Now().UTC()
	templateRepo := repo.NewTemmplateRepository(template.GroupCode)
	if templateRepo == nil {
		tempChannel <- common.MSG_UNATHORIZED_ACCESS
	}
	var err error
	newTemplate := !utility.IsPrimaryIDValid(template.TemplateID)
	if newTemplate {
		template.CreateDate = template.LastModifiedDate
		err = templateRepo.Save(template)
	} else {
		err = templateRepo.Update(template)
	}
	if err != nil {
		fmt.Println(err)
		errorMsg := utility.GetErrorMsg(err)
		if len(errorMsg) > 0 {
			tempChannel <- errorMsg
		}
		tempChannel <- common.MSG_REQUEST_FAILED
	}
	tempChannel <- common.MSG_SAVE_SUCCESS
}

func saveTestPaper(testPaper *model.TestPaper, saveAsDraft bool,
	wg *sync.WaitGroup, paperChannel chan<- string) {
	wg.Add(1)
	defer cleanPanic(paperChannel, wg)
	if testPaper == nil {
		paperChannel <- ""
	}
	testPaper.LastModifiedDate = time.Now().UTC()
	if saveAsDraft {
		testPaper.Status = enums.CurrentTestStatus.DRAFT
	} else {
		testPaper.Status = enums.CurrentTestStatus.RELEASE
	}
	testPaperRepo := repo.NewTestPaperRepository(testPaper.GroupCode)
	if testPaperRepo == nil {
		paperChannel <- common.MSG_UNATHORIZED_ACCESS
	}
	var err error
	newTestPaper := !utility.IsPrimaryIDValid(testPaper.TestID)
	if newTestPaper {
		testPaper.CreateDate = testPaper.LastModifiedDate
		err = testPaperRepo.Save(testPaper)
	} else {
		err = testPaperRepo.Update(testPaper)
	}
	if err != nil {
		fmt.Println(err)
		errorMsg := utility.GetErrorMsg(err)
		if len(errorMsg) > 0 {
			paperChannel <- errorMsg
		}
		paperChannel <- common.MSG_REQUEST_FAILED
	}
	paperChannel <- common.MSG_SAVE_SUCCESS
}

func cleanPanic(channel chan<- string, wg *sync.WaitGroup) {
	wg.Done()
	if rec := recover(); rec != nil {
		channel <- common.MSG_REQUEST_FAILED
	}
}

//FetchAllDrafts gets all drafted test papers under a teacher
func FetchAllDrafts(groupCode, teacherID string) *dto.AppResponseDto {
	if utility.IsStringIDValid(teacherID) {
		testPaperRepo := repo.NewTestPaperRepository(groupCode)
		if testPaperRepo == nil {
			return utility.GetErrorResponse(common.MSG_UNATHORIZED_ACCESS)
		}
		drafts, err := testPaperRepo.FindAll(bson.ObjectIdHex(teacherID))
		if err != nil {
			fmt.Println(err.Error())
			errorMsg := utility.GetErrorMsg(err)
			if len(errorMsg) > 0 {
				return utility.GetErrorResponse(errorMsg)
			}
			return utility.GetErrorResponse(common.MSG_REQUEST_FAILED)
		}
		return utility.GetSuccessResponse(drafts)
	}
	return utility.GetErrorResponse(common.MSG_INVALID_ID)
}

//FetchSinglePaper gets one drafted paper as per testId
func FetchSinglePaper(groupCode, testPaperID string) *dto.AppResponseDto {
	if utility.IsStringIDValid(testPaperID) {
		testPaperRepo := repo.NewTestPaperRepository(groupCode)
		if testPaperRepo == nil {
			return utility.GetErrorResponse(common.MSG_UNATHORIZED_ACCESS)
		}
		testPaper, err := testPaperRepo.FindOne(bson.ObjectIdHex(testPaperID))
		if err != nil {
			fmt.Println(err.Error())
			errorMsg := utility.GetErrorMsg(err)
			if len(errorMsg) > 0 {
				return utility.GetErrorResponse(errorMsg)
			}
			return utility.GetErrorResponse(common.MSG_REQUEST_FAILED)
		}
		return utility.GetSuccessResponse(testPaper)
	}
	return utility.GetErrorResponse(common.MSG_INVALID_ID)
}

//FetchAllTemplates gets all templates under a teacher
func FetchAllTemplates(groupCode, teacherID string) *dto.AppResponseDto {
	if utility.IsStringIDValid(teacherID) {
		templateRepo := repo.NewTemmplateRepository(groupCode)
		if templateRepo == nil {
			return utility.GetErrorResponse(common.MSG_UNATHORIZED_ACCESS)
		}
		templates, err := templateRepo.FindAll(bson.ObjectIdHex(teacherID))
		if err != nil {
			fmt.Println(err.Error())
			errorMsg := utility.GetErrorMsg(err)
			if len(errorMsg) > 0 {
				return utility.GetErrorResponse(errorMsg)
			}
			return utility.GetErrorResponse(common.MSG_REQUEST_FAILED)
		}
		return utility.GetSuccessResponse(templates)
	}
	return utility.GetErrorResponse(common.MSG_INVALID_ID)
}

//FetchSingleTemplate gets one template as per templateId
func FetchSingleTemplate(groupCode, testPaperID string) *dto.AppResponseDto {
	if utility.IsStringIDValid(testPaperID) {
		templateRepo := repo.NewTemmplateRepository(groupCode)
		if templateRepo == nil {
			return utility.GetErrorResponse(common.MSG_UNATHORIZED_ACCESS)
		}
		template, err := templateRepo.FindOne(bson.ObjectIdHex(testPaperID))
		if err != nil {
			fmt.Println(err.Error())
			errorMsg := utility.GetErrorMsg(err)
			if len(errorMsg) > 0 {
				return utility.GetErrorResponse(errorMsg)
			}
			return utility.GetErrorResponse(common.MSG_REQUEST_FAILED)
		}
		return utility.GetSuccessResponse(template)
	}
	return utility.GetErrorResponse(common.MSG_INVALID_ID)
}
