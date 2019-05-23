package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const (
	demoLabelKey = "com.infosiftr.dockercon-demo.active"
	demoLabelVal = "yes"

	randomBytes = 1024 * 1024
)

func kubectl(args ...string) (string, error) {
	fmt.Fprintf(os.Stderr, "$ kubectl %q\n", args)
	cmd := exec.Command("kubectl", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	outStr := strings.TrimSpace(out.String())
	fmt.Fprintf(os.Stderr, "%s\n\n", outStr)
	return outStr, err
}

func js(in interface{}) string {
	b, err := json.MarshalIndent(in, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func apiNodes(w http.ResponseWriter, r *http.Request) {
	nodeServices := map[string]map[string]string{} // [node][service] = CurrentState
	servicesTxt, err := kubectl("get", "pods", "-l", "app=blinky-server", "-o", `go-template={{ range .items }}{{ .spec.nodeName }}|{{ .metadata.name }}|{{ .status.phase }}{{ range .status.containerStatuses }}{{ range .state }}{{ if .reason }} / {{ .reason }}{{ end }}{{ end }}{{ end }}{{ "\n" }}{{ end }}`)
	if err != nil {
		log.Print("kubectl get pods: ", err)
		return
	}
	for _, service := range strings.Split(servicesTxt, "\n") {
		if service == "" {
			// ignore empty lines (likely means everything is turned off)
			continue
		}
		serviceParts := strings.SplitN(service, "|", 3)
		if _, ok := nodeServices[serviceParts[0]]; !ok {
			nodeServices[serviceParts[0]] = map[string]string{}
		}
		nodeServices[serviceParts[0]][serviceParts[1]] = serviceParts[2]
	}

	nodesTxt, err := kubectl("get", "nodes", "-l", "!node-role.kubernetes.io/master", "-o", `go-template={{ range .items }}{{ .metadata.name }}|{{ .status.nodeInfo.operatingSystem }}|{{ .status.nodeInfo.architecture }}|{{ .metadata.labels.app }}{{ "\n" }}{{ end }}`)
	if err != nil {
		log.Print("kubectl get nodes: ", err)
		return
	}
	ret := []map[string]interface{}{}
	nodes := strings.Split(nodesTxt, "\n")
	sort.Strings(nodes)
	for _, node := range nodes {
		nodeParts := strings.SplitN(node, "|", 4)
		nodeRet := map[string]interface{}{
			"ID":       nodeParts[0],
			"Hostname": nodeParts[0],

			"OS": nodeParts[1],
			"Architecture": nodeParts[2],
		}

		nodeRet["Services"] = nodeServices[nodeRet["Hostname"].(string)]

		nodeRet["DemoActive"] = nodeParts[3] == "blinky-node"

		ret = append(ret, nodeRet)
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, js(ret))
}

func apiNodeActivate(w http.ResponseWriter, r *http.Request) {
	node := r.URL.Query().Get("node")
	if node == "" {
		return
	}
	_, err := kubectl("label", "node", node, "app=blinky-node", "--overwrite")
	if err != nil {
		log.Print("kubectl node update: ", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, js(true))
}

func apiNodeDeactivate(w http.ResponseWriter, r *http.Request) {
	node := r.URL.Query().Get("node")
	if node == "" {
		return
	}
	_, err := kubectl("label", "node", node, "app=blinky-nope", "--overwrite")
	if err != nil {
		log.Print("kubectl node update: ", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, js(true))
}

func apiEcho(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && r.Body != nil {
		w.Header().Add("Content-Type", contentType)
		io.Copy(w, r.Body)
	} else {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, js(true))
	}
}

func wwwHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/", http.StatusFound)
}

func www() {
	http.HandleFunc("/api/nodes", apiNodes)
	http.HandleFunc("/api/node/activate", apiNodeActivate)
	http.HandleFunc("/api/node/deactivate", apiNodeDeactivate)
	http.HandleFunc("/api/echo", apiEcho)

	http.HandleFunc("/", wwwHome)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func md5sum(r io.Reader) string {
	hash := md5.New()
	io.Copy(hash, r)
	return hex.EncodeToString(hash.Sum(nil))
}

func workerApiDump(apiServer string) {
	echoEndpoint := fmt.Sprintf("http://%s:8080/api/echo", apiServer)

	buf := make([]byte, randomBytes)
	_, err := rand.Read(buf)
	if err != nil {
		log.Print("rand error: ", err)
		return
	}

	log.Print("sending  ", md5sum(bytes.NewReader(buf)))

	resp, err := http.Post(echoEndpoint, "application/octet-stream", bytes.NewReader(buf))
	if err != nil {
		log.Print("post error: ", err)
		return
	}
	defer resp.Body.Close()

	log.Print("received ", md5sum(resp.Body))
}

// usage: dockercon-demo [api-server]
//
// run without arguments to start an api-server instance
// run with a single argument to blast data at "/api/echo" endpoint of the given api-server instance (blink blink blink go the lights)
func main() {
	if len(os.Args) == 1 {
		www()
	} else {
		if len(os.Args) != 2 {
			log.Fatalf("wrong number of arguments! expected 1 or 2, not %d", len(os.Args)-1)
		}
		apiServer := os.Args[1]
		for {
			workerApiDump(apiServer)
		}
	}
}
