package main

import (
    "net/http"
    "net"
    "log"
    "encoding/json"
    "os"
    "strings"
    "errors"
)

type TestHandler struct {
}

// getIP returns request real ip.
func getIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	return "", errors.New("no valid ip found")
}

func Healthz(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(http.StatusOK)
    log.Printf("healthz")
    w.Write([]byte(string("Healthz")))
}
//ServeHTTP方法，绑定TestHandler
func (th *TestHandler)ServeHTTP(w http.ResponseWriter, r *http.Request){
    reqHeader := r.Header
    version := os.Getenv("VERSION")
    reqHeader["ENV_VERSION"] = []string{version}
    for k,v := range reqHeader {
        log.Printf("%v:%v", k, v)
        w.Header().Add(k, strings.Join(v, ";"))
    }
    ip, err := getIP(r)
    if err != nil {
        log.Printf("%s", err)
    } else {
        log.Printf("ip:%v", ip)
    }
    data,_ := json.Marshal(reqHeader)
    w.Write(data)
}

func main(){
    http.Handle("/", &TestHandler{})
    http.HandleFunc("/healthz", Healthz)
    http.ListenAndServe("127.0.0.1:8000",nil)
}