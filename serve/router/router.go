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

func Run(port int) {
	gin.SetMode(gin.DebugMode)
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

	api.POST("/project_list/save", controllers.SaveProjectList)
	api.POST("/choose_team", controllers.ChooseTeam)
	//api.POST("/issue_type_list", controllers.IssueTypeList)
	api.POST("/issue_type_list/save", controllers.SaveIssueTypeList)

	api.GET("/import/reset", controllers.Reset)
	api.POST("/import/start", controllers.StartImport)
	api.POST("/import/pause", controllers.PauseImport)
	api.POST("/import/continue", controllers.ContinueImport)
	//api.POST("/import/stop", controllers.StopImport)
	//api.GET("/import/progress", controllers.ImportProgress)
	api.GET("/import/log", controllers.GetAllImportLog)
	api.GET("/import/log/start_line/:start_line", controllers.GetImportLog)
	api.GET("/import/log/download/all", controllers.DownloadLogFile)
	api.GET("/import/log/download/current", controllers.DownloadCurrentLogFile)
	api.GET("/import/scope", controllers.GetScope)

	api.GET("/common_config", controllers.Config)
	api.POST("/login", controllers.Login)
	api.POST("/logout", middlewares.CheckLogin, controllers.Logout)
	api.GET("/jira_config", middlewares.CheckLogin, controllers.UserJiraConfig)
	api.POST("/check_jira_path_exist", middlewares.CheckLogin, controllers.CheckJiraPathExist)
	api.POST("/jira_backup_list", middlewares.CheckLogin, controllers.JiraBackUpList)
	api.POST("/resolve/start", middlewares.CheckLogin, controllers.StartResolve)
	api.GET("/resolve/progress", middlewares.CheckLogin, controllers.ResolveProgress)
	api.POST("/resolve/stop", middlewares.CheckLogin, controllers.StopResolve)
	api.GET("/resolve/result", middlewares.CheckLogin, controllers.ResolveResult)
	api.GET("/team_list", middlewares.CheckLogin, controllers.TeamList)
	api.GET("/history_config/project", middlewares.CheckLogin, controllers.ProjectHistoryConfig)
	api.GET("/project_list", middlewares.CheckLogin, controllers.ProjectList)
	api.POST("/history_config/project", middlewares.CheckLogin, controllers.SetProjectHistoryConfig)
	api.POST("/check_disk", middlewares.CheckLogin, controllers.CheckProjectDisk)
	//api.POST("/issue_type_list", controllers.IssueTypeList)

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
