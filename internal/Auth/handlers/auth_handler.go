package handlers

type AuthHandler struct {
}

// func (handler *CustomerHandler) Login(ctx *gin.Context) {
// 	var reqData model.LoginReq

// 	if err := ctx.ShouldBindJSON(&reqData); err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid email or password"))
// 		return
// 	}

// 	data, err := handler.Usecase.FindByEmail(ctx, reqData.Email)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
// 		return
// 	}
// 	if !template.CheckPasswordHash(reqData.Password, data.Password) {
// 		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid email or password"))
// 		return
// 	}
// 	// cache token
// 	// _, err = auth.NewTokenGenerator(db.Redis(), *data)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, response.Success("Login Success"))
// }

// func (handler *CustomerHandler) Create(ctx *gin.Context) {
// 	var data model.CreateReq

// 	if err := ctx.ShouldBindJSON(&data); err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.BadRequest("Invalid input data"))
// 		return
// 	}
// 	_, err := handler.Usecase.Create(ctx, data)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, response.Success("Success Create User"))
// }
