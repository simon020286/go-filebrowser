package main

import (
	"encoding/json"
	"filebrowser/models"
	"filebrowser/utils"
	"filebrowser/ws"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/buaazp/fasthttprouter"
	"github.com/chebyrash/promise"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

var (
	fu    utils.IFileUtils
	tasks utils.Tasks
	hub   *ws.Hub
)

func main() {
	basePath := os.Getenv("BASEPATH")
	fu = utils.NewFileUtils(basePath)
	tasks = utils.NewTasks()
	hub = ws.NewHub()
	go hub.Run()

	router := fasthttprouter.New()

	router.ServeFiles("/react-filemanager/*filepath", "static")
	router.GET("/filemanager/list", list)
	router.POST("/filemanager/dir/create", create)
	router.POST("/filemanager/items/copy", copyMoveHandler(true))
	router.POST("/filemanager/items/move", copyMoveHandler(false))
	router.POST("/filemanager/items/remove", remove)
	router.POST("/filemanager/items/upload", upload)
	router.GET("/filemanager/file/content", content)
	router.GET("/filemanager/tasks", taskList)
	router.GET("/filemanager/ws", provaWs)
	router.GET("/ws", hub.ServeWs)

	tasks.OnAddedTask = func(id uuid.UUID, task *models.Task) {
		hub.Broadcast("task added", (*task))
	}

	tasks.OnEndedTask = func(id uuid.UUID, task *models.Task) {
		hub.Broadcast("task ended", (*task))
	}

	log.Fatal(fasthttp.ListenAndServe(":8080", cors(router.Handler)))
}

func provaWs(ctx *fasthttp.RequestCtx) {
	hub.Broadcast("prova", models.Task{Name: "Prova"})
}

func list(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	src := string(args.Peek("path"))
	src = path.Clean(src)
	files, err := fu.GetFiles(src)
	if err != nil {
		models.NewAPIResponse(models.OptionError("Cannot read folder", err)).Send(ctx)
	} else {
		models.NewAPIResponse(models.OptionData(files)).Send(ctx)
	}
}

func upload(ctx *fasthttp.RequestCtx) {
	multipart, err := ctx.MultipartForm()
	if err != nil {
		models.NewAPIResponse(models.OptionError("Cannot read folder", err)).Send(ctx)
		return
	}

	var promises []*promise.Promise
	files := multipart.File["file[]"]

	for index := 0; index < len(files); index++ {
		file, err := files[index].Open()
		if err != nil {
			models.NewAPIResponse(models.OptionError("An error ocurred uploading files", err)).Send(ctx)
			return
		}
		defer file.Close()
		promises = append(promises, fu.SavePromise(file, files[index].Filename, multipart.Value["path"][0]))
	}

	promise.All(promises...).
		Then(func(data interface{}) interface{} {
			models.NewAPIResponse(models.OptionData(data)).Send(ctx)
			return data
		}).
		Catch(func(err error) error {
			models.NewAPIResponse(models.OptionError("An error ocurred uploading files", err)).Send(ctx)
			return err
		}).Await()
}

func create(ctx *fasthttp.RequestCtx) {
	var body struct {
		Path      string `json:"path"`
		Directory string `json:"directory"`
	}

	json.Unmarshal(ctx.PostBody(), &body)
	err := fu.CreateFolder(path.Clean(path.Join(body.Path, body.Directory)))
	if err != nil {
		models.NewAPIResponse(models.OptionError("Unknown error creating folder", err)).Send(ctx)
	}
	models.NewAPIResponse(models.OptionData(nil)).Send(ctx)
}

func copyMoveHandler(copy bool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var body struct {
			Path        string   `json:"path"`
			Filenames   []string `json:"filenames"`
			Destination string   `json:"destination"`
		}
		json.Unmarshal(ctx.PostBody(), &body)

		var promises []*promise.Promise

		for _, filename := range body.Filenames {
			promises = append(promises, promise.New(func(resolve func(interface{}), reject func(error)) {
				var (
					size int64
					err  error
					task models.Task
				)
				src := path.Join(body.Path, filename)
				dst := path.Join(body.Destination, filename)

				if copy {
					task = models.NewCopyTask(src, dst)
					tasks.Add(&task)
					size, err = fu.Copy(src, dst)
				} else {
					task = models.NewMoveTask(src, dst)
					tasks.Add(&task)
					size, err = fu.Move(src, dst)
				}
				if err != nil {
					task.End(err)
					reject(err)
					return
				}
				task.End(nil)
				resolve(size)
				return
			}))
		}

		promise.All(promises...).
			Then(func(data interface{}) interface{} {
				models.NewAPIResponse(models.OptionData(data)).Send(ctx)
				return data
			}).
			Catch(func(err error) error {
				models.NewAPIResponse(models.OptionError("An error ocurred copying files", err)).Send(ctx)
				return err
			}).Await()
	}
}

func remove(ctx *fasthttp.RequestCtx) {
	var body struct {
		Path      string   `json:"path"`
		Filenames []string `json:"filenames"`
		Recursive bool     `json:"recursive"`
	}
	json.Unmarshal(ctx.PostBody(), &body)

	var promises []*promise.Promise

	for _, filename := range body.Filenames {
		promises = append(promises, fu.DeletePromise(path.Join(body.Path, filename)))
	}

	promise.All(promises...).
		Then(func(data interface{}) interface{} {
			models.NewAPIResponse(models.OptionData(data)).Send(ctx)
			return data
		}).
		Catch(func(err error) error {
			models.NewAPIResponse(models.OptionError("An error ocurred copying files", err)).Send(ctx)
			return err
		}).Await()
}

func content(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	src := string(args.Peek("path"))
	src = path.Clean(src)
	filename := filepath.Base(src)
	ctx.Response.Header.Add("Content-Disposition", "attachment; filename=\""+filename+"\"")
	ctx.SendFile(fu.GetAbsolutePath(src))
}

func taskList(ctx *fasthttp.RequestCtx) {
	list := tasks.All()
	models.NewAPIResponse(models.OptionData(list)).Send(ctx)
	tasks.CleanEnded()
}

func cors(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		ctx.Response.Header.Add("Access-Control-Allow-Headers", "X-Requested-With,content-type,path")
		ctx.Response.Header.Add("Access-Control-Expose-Headers", "Content-Disposition")
		ctx.Response.Header.Add("Access-Control-Allow-Credentials", "true")

		next(ctx)
	}
}
