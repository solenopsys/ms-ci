package main

import (
	"alexstorm-hsm-ci/internal/kube"
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	bl_kubernetes_tools "github.com/solenopsys/bl-kubernetes-tools"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

var clientSet *kubernetes.Clientset
var c client.Client
var devMode = os.Getenv("developerMode") == "true"

const ConfigmapName = "helm-repositories"
const NameSpace = "default"

type Execution struct {
	id        int
	params    string
	container string
}

func CheckError(err error) {
	if err != nil {

		klog.Error(err)
	}
}

func list(db *sql.DB) {
	rows, err := db.Query("select * from alexstorm_shockwaves.execution;")
	CheckError(err)

	p := Execution{}

	for rows.Next() {
		err = rows.Scan(&p.id, &p.params, &p.container)
		CheckError(err)
		println("ROW")
		println(p.params)
		println(p.id)
	}

}

var Mode string

const DEV_MODE = "dev"

func init() {
	flag.StringVar(&Mode, "mode", "", "a string var")
}

func main() {
	flag.Parse()
	devMode := Mode == DEV_MODE

	if devMode {
		err := godotenv.Load("configs/local.env")
		CheckError(err)
	}

	host := os.Getenv("postgres.Host")
	port, err := strconv.ParseInt(os.Getenv("postgres.Port"), 10, 16)
	CheckError(err)
	user := os.Getenv("postgres.User")
	password := os.Getenv("postgres.Password")
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, "Converged")

	// open database
	db, err := sql.Open("postgres", psqlconn)

	CheckError(err)

	list(db)
	clientSet, c = bl_kubernetes_tools.CreateKubeConfig(devMode)

	gitRepoName := "solenopsys-hsm-ci"
	gitHost := "git.solenopsys.org"

	kube.CreateJobFunc(clientSet,
		gitRepoName,
		"ci-build-job7",
		gitHost,
		"linux/amd64,linux/arm64",
		"/workspace/"+gitRepoName+"/cic/jobs/test",
		"registry.solenopsys.org/"+gitRepoName,
		map[string]string{
			"REPO_NAME":    gitRepoName,
			"USER_INFO":    "admin:root@",
			"GIT_HOST":     gitHost,
			"GO_MAIN_FILE": "/sources/cmd/app/main.go",
		},
	)
	//	template := zmq_connector.HsTemplate{Pf: processingFunction()}
	//	template.Init()
}
