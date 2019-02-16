package approuter

import (
	"github.com/gin-gonic/gin"

	"intelliq/app/controller"
)

var mrouter *gin.Engine

//AddRouters adding routes
func AddRouters(router *gin.Engine) {
	mrouter = router
	addMetaRouters()
	addUserRouters()
	addSchoolRouters()
	addGroupRouters()
	addQuestionRouters()
}

func addMetaRouters() {
	metaRoutes := mrouter.Group("/meta")
	{
		metaRoutes.POST("/add", controller.AddMetaData)
		metaRoutes.PUT("/update", controller.UpdateMetaData)
		metaRoutes.GET("/read", controller.ReadMetaData)
	}
}

func addUserRouters() {
	userRoutes := mrouter.Group("/user")
	{
		userRoutes.POST("/add", controller.AddNewUser)
		userRoutes.PUT("/update", controller.UpdateUserProfile)
		userRoutes.GET("/all/admins/:groupId", controller.ListAllSchoolAdmins)
		userRoutes.GET("/all/school/:schoolId", controller.ListAllTeachers)
		userRoutes.GET("/all/school/:schoolId/:roleType", controller.ListSelectedTeachers)
		userRoutes.PUT("/role/transfer/:roleType/:fromUser/:toUser", controller.TransferRole)
		userRoutes.DELETE("/remove/:schoolId/:userId", controller.RemoveUserFromSchool)
		userRoutes.POST("/bulk/add", controller.AddBulkUsers)
		userRoutes.POST("/bulk/update", controller.UpdateBulkUsers)
		userRoutes.POST("/login", controller.AuthenticateUser)
		userRoutes.GET("/logout/:userId", controller.Logout)
		userRoutes.GET("/info/:key/:val", controller.ListUserByMobileOrID)
	}
}

func addSchoolRouters() {
	schoolRoutes := mrouter.Group("/school")
	{
		schoolRoutes.POST("/add", controller.AddNewSchool)
		schoolRoutes.GET("/all/:key/:val", controller.ListAllSchools)
		schoolRoutes.PUT("/update", controller.UpdateSchoolProfile)
		schoolRoutes.GET("info/:key/:val", controller.ListSchoolByCodeOrID)
	}
}

func addGroupRouters() {
	groupRoutes := mrouter.Group("/group")
	{
		groupRoutes.POST("/add", controller.AddNewGroup)
		groupRoutes.PUT("/update", controller.UpdateGroup)
		groupRoutes.GET("/all/:restrict", controller.ListAllGroups)
		groupRoutes.GET("info/:key/:val", controller.ListGroupByCodeOrID)
	}
}

func addQuestionRouters() {
	//quesRoutes := mrouter.Group("/question"){}

}
