package main

import (
	"alexstorm-hsm-ci/internal/kube"
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
		godotenv.Load("configs/local.env")
	}

	host := os.Getenv("postgres.Host")
	port, err := strconv.ParseInt(os.Getenv("postgres.Port"), 10, 16)
	CheckError(err)
	user := os.Getenv("postgres.User")
	password := os.Getenv("postgres.Password")
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, "softconverged")

	// open database
	db, err := sql.Open("postgres", psqlconn)

	CheckError(err)

	list(db)
	clientSet, c = kube.СreateKubeConfig(devMode)
	name := "solenopsys-hsm-ci"
	kube.CreateJobFunc(clientSet,
		name,
		"ci-build-job29",
		"git.alexstorm.solenopsys.org",
		"linux/amd64,linux/arm64",
		"/workspace/"+name+"/cic/jobs/test",
		"registry.alexstorm.solenopsys.org/"+name,
	)
	//	template := zmq_connector.HsTemplate{Pf: processingFunction()}
	//	template.Init()
}
