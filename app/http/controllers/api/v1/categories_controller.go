package v1

import (
	"gohub/app/models/category"
	"gohub/app/requests"
	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
)

type CategoriesController struct {
	BaseAPIController
}

// func (ctrl *CategoriesController) Index(c *gin.Context) {
// 	categories := category.All()
// 	response.Data(c, categories)
// }

// func (ctrl *CategoriesController) Show(c *gin.Context) {
// 	categoryModel := category.Get(c.Param("id"))
// 	if categoryModel.ID == 0 {
// 		response.Abort404(c)
// 		return
// 	}
// 	response.Data(c, categoryModel)
// }

func (ctrl *CategoriesController) Store(c *gin.Context) {

	request := requests.CategoryRequest{}
	if ok := requests.Validate(c, &request, requests.CategorySave); !ok {
		return
	}

	categoryModel := category.Category{
		Name:        request.Name,
		Description: request.Description,
	}
	categoryModel.Create()
	if categoryModel.ID > 0 {
		response.Created(c, categoryModel)
	} else {
		response.Abort500(c, "创建失败, 请稍后尝试~")
	}
}

func (ctrl *CategoriesController) Update(c *gin.Context) {

	categoryModel := category.Get(c.Param("id"))
	if categoryModel.ID == 0 {
		response.Abort404(c)
		return
	}

	request := requests.CategoryRequest{}
	if ok := requests.Validate(c, &request, requests.CategorySave); !ok {
		return
	}

	categoryModel.Name = request.Name
	categoryModel.Description = request.Description
	rowsAffected := categoryModel.Save()

	if rowsAffected > 0 {
		response.Data(c, categoryModel)
	} else {
		response.Abort500(c, "更新失败, 请稍后尝试~")
	}
}

// func (ctrl *CategoriesController) Delete(c *gin.Context) {

// 	categoryModel := category.Get(c.Param("id"))
// 	if categoryModel.ID == 0 {
// 		response.Abort404(c)
// 		return
// 	}

// 	if ok := policies.CanModifyCategory(c, categoryModel); !ok {
// 		response.Abort403(c)
// 		return
// 	}

// 	rowsAffected := categoryModel.Delete()
// 	if rowsAffected > 0 {
// 		response.Success(c)
// 		return
// 	}

// 	response.Abort500(c, "删除失败，请稍后尝试~")
// }
