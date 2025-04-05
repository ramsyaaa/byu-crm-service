package routes

import (
	"byu-crm-service/modules/detail-community-member/http"
	"byu-crm-service/modules/detail-community-member/repository"
	"byu-crm-service/modules/detail-community-member/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DetailCommunityMemberRouter(app *fiber.App, db *gorm.DB) {
	detailCommunityMemberRepo := repository.NewDetailCommunityMemberRepository(db)
	detailCommunityMemberService := service.NewDetailCommunityMemberService(detailCommunityMemberRepo)
	detailCommunityMemberHandler := http.NewDetailCommunityMemberHandler(detailCommunityMemberService)

	http.DetailCommunityMemberRoute(app, detailCommunityMemberHandler)

}
