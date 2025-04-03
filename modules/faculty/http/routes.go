package http

import "github.com/gofiber/fiber/v2"

func FacultyRoutes(router fiber.Router, handler *FacultyHandler) {
	router.Get("/faculties", handler.GetAllFaculties)
	router.Get("/faculties/:id", handler.GetFacultyByID)
	router.Post("/faculties", handler.CreateFaculty)
	router.Put("/faculties/:id", handler.UpdateFaculty)
}
