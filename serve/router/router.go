package router

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"

	"github.com/bangwork/import-tools/serve/controllers"
	"github.com/bangwork/import-tools/serve/middlewares"
)

//go:embed dist lang
var FS embed.FS
var FSWithLogin embed.FS

func Run(port int) {
	gin.SetMode(gin.ReleaseMode)
	api := gin.Default()
	api.Use(middlewares.Recovery(), middlewares.Logger())
	api.Use(GinI18nLocalize())
	api.Use(middlewares.Cors())

	temple := template.Must(template.New("").ParseFS(FS, "dist/index.html"))
	api.SetHTMLTemplate(temple)

	fe, err := fs.Sub(FS, "dist/assets")
	if err != nil {
		log.Println("embed dist assets err", err)
		return
	}
	api.StaticFS("/assets", http.FS(fe))

	fe, _ = fs.Sub(FS, "dist/public")
	api.StaticFS("/public", http.FS(fe))

	api.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	api.POST("/check_path_exist", controllers.CheckPathExist)
	api.POST("/check_jira_path_exist", controllers.CheckJiraPathExist)
	api.POST("/jira_backup_list", controllers.JiraBackUpList)
	api.POST("/resolve/start", controllers.StartResolve)
	api.GET("/resolve/progress", controllers.ResolveProgress)
	api.POST("/resolve/stop", controllers.StopResolve)
	api.GET("/resolve/result", controllers.ResolveResult)
	api.GET("/project_list", controllers.ProjectList)
	api.POST("/project_list/save", controllers.SaveProjectList)
	api.POST("/choose_team", controllers.ChooseTeam)
	api.POST("/issue_type_list", controllers.IssueTypeList)
	api.POST("/issue_type_list/save", controllers.SaveIssueTypeList)

	api.GET("/import/reset", controllers.Reset)
	api.POST("/import/start", controllers.StartImport)
	api.POST("/import/pause", controllers.PauseImport)
	api.POST("/import/continue", controllers.ContinueImport)
	api.POST("/import/stop", controllers.StopImport)
	api.GET("/import/progress", controllers.ImportProgress)
	api.GET("/import/log", controllers.GetAllImportLog)
	api.GET("/import/log/start_line/:start_line", controllers.GetImportLog)
	api.GET("/import/log/download/all", controllers.DownloadLogFile)
	api.GET("/import/log/download/current", controllers.DownloadCurrentLogFile)
	api.GET("/import/scope", controllers.GetScope)

	//apiLogin := gin.Default()
	//apiLogin.Use(middlewares.Recovery(), middlewares.Logger())
	//apiLogin.Use(GinI18nLocalize())
	//apiLogin.Use(middlewares.Cors())
	//
	//temple1 := template.Must(template.New("").ParseFS(FS, "dist/index.html"))
	//apiLogin.SetHTMLTemplate(temple1)
	//
	//fe1, err1 := fs.Sub(FS, "dist/assets")
	//if err1 != nil {
	//	log.Println("embed dist assets err", err1)
	//	return
	//}
	//apiLogin.StaticFS("/assets", http.FS(fe1))
	//
	//fe1, _ = fs.Sub(FSWithLogin, "dist/public")
	//apiLogin.StaticFS("/public", http.FS(fe1))
	//
	//apiLogin.GET("/", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "index.html", gin.H{})
	//})

	api.POST("/login", controllers.Login)
	api.POST("/logout", middlewares.CheckLogin, controllers.Logout)

	api.Run(fmt.Sprintf(":%d", port))
}

func GinI18nLocalize() gin.HandlerFunc {
	return i18n.Localize(
		i18n.WithBundle(&i18n.BundleCfg{
			RootPath:         "./lang",
			AcceptLanguage:   []language.Tag{language.Chinese, language.English},
			DefaultLanguage:  language.Chinese,
			FormatBundleFile: "toml",
			UnmarshalFunc:    toml.Unmarshal,
			Loader:           &i18n.EmbedLoader{FS: FS},
		}),
	)
}
