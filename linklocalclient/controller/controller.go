package controller

import(
	"github.com/ccwings/log"
	"gopkg.in/macaron.v1"
	"github.com/martini-contrib/cors"
	"linklocal/auth"
)

func InitRouter(m *macaron.Macaron){
	initBaseRouter(m)
	initAuthRouter(m)
}

func initBaseRouter(m *macaron.Macaron){
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:[]string{"*"},
		AllowMethods:[]string{"GET","POST","PUT","DELETE"},
		AllowHeaders:[]string{"Origin", "x-requested-with", "Content-Type", "Content-Range", "Content-Disposition", "Content-Description", "token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	log.Info("BaseRouter Init func starting")
	m.Get("/",authorize,info)
}

func initAuthRouter(m *macaron.Macaron){
	m.Post("/login",Login)
	m.Post("/logout",Logout)
	m.Post("/register",Register)
	m.Group("/auth",func(){
		m.Post("/checkusername",CheckUsername)
		m.Post("/checkemail",CheckEmail)
		m.Group("/users",func(){
			m.Post("",authorize,CreateUser)
			m.Delete("/:id",authorize,DeleteUser)
			m.Put("/:id",authorize,UpdateUser)
			m.Get("",authorize,ListUsers)
			m.Get("/:id",authorize,GetUser)
			m.Post("/:id/changepassword",authorize,ChangeUserPassword)
			m.Post("/:id/changeendtime",authorize,ChangeUserEndTime)
			//user_group
			m.Get("/:id/groups",authorize,GetUserBindingGroups)
		})
		m.Group("/groups",func(){
			m.Post("",authorize,CreateGroup)
			m.Delete("/:id",authorize,DeleteGroup)
			m.Put("/:id",authorize,UpdateGroup)
			m.Get("",authorize,ListGroups)
			m.Get("/:id",authorize,GetGroup)
			m.Post("/:id/changeendtime",authorize,ChangeGroupEndTime)
			m.Get("/:id/parent",authorize,GetGroupParent)
			m.Get("/:id/children",authorize,GetGroupChildren)
			//user_group
			m.Get("/:id/users",authorize,GetGroupBindingUsers)
			m.Post("/:id/binduser",authorize,BindUserToGroup)
			m.Post("/:id/unbinduser",authorize,UnbindUserToGroup)
			m.Post("/:id/setusertype",authorize,SetUserTypeForGroup)
			m.Post("/:id/setuserpriority",authorize,SetUserPriorityForGroup)
		})
	})
}


func info(ctx *macaron.Context){
	ctx.JSON(200, "LinkLocal Ver 1.0.0, API Ver 1.0.0")
	return
}

func authorize(ctx *macaron.Context) {
	tokenStr := ctx.Req.Header.Get("token")
	loginInfo := make(map[string]interface{})
	if tokenStr == "" {
		loginInfo["error"] = "Action Need Login"
		ctx.JSON(401, loginInfo)
		return
	}
	if !auth.TokenAuthorization(tokenStr) {
		loginInfo["error"] = "Token Expired or Not Authorized"
		ctx.JSON(401, loginInfo)
		return
	}
}