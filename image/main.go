package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var tpl = template.Must(template.ParseFiles("index.html"))
var root = "/tmp/kuberesults/kube/"

//Result struct is
type Result struct {
	Result string `json: "result"`
	Desc   string `json: "desc"`
}

//File struct is
type File struct {
	Name string `json: "name"`
}

func kubeCall() {
	url := "https://console-int.okd.local:443/apis/extensions/v1beta1/namespaces/kube-bench/daemonsets"
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLWJlbmNoIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImdvbGFuZy1zYS10b2tlbi1kYjZ4bSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJnb2xhbmctc2EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIxZjZlMjE1YS01NDE5LTExZWItOGM0MC0wMDBjMjk0M2FjNzYiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1iZW5jaDpnb2xhbmctc2EifQ.o4h_cEQzLksXCq58J2EdDogj95qjfxQjI9mNDyXFP1GJhTwbyM8CQQmiT77ic0vb4SRGjHSh_rZyn_FkbFi4NiUAd_K4aXmlv6EOyh5moBlauw_2g7ezBx_MJ-QZVdFfbvsyqICG70KlCBnbe6ARu3bHquWRPSQKx0wvCJpjwinUJHPZrlbhO7VIezAG4xHjjUdwH2v-sduw0z6OARX4oeV4XGZso6L2-O1f_1Dx67Lx22X9cLNe3_Iox14w44EDTpkbiViruZHrgwJ3zoehPIDNO_ZOR-40aCP1ineSqgLkxQQ9re4_kLo4AmT8Gv1M5CahOhDPBPEK2plBFR-C4A"
	var bearer = "Bearer " + token
	var jsonStr = []byte(`{
		"apiVersion": "extensions/v1beta1",
		"kind": "DaemonSet",
		"metadata": {
		   "name": "kube-bench"
		},
		"spec": {
		   "selector": {
			  "matchLabels": {
				 "app": "kube-bench"
			  }
		   },
		   "template": {
			  "metadata": {
				 "labels": {
					"app": "kube-bench"
				 }
			  },
			  "spec": {
				 "hostPID": true,
				 "containers": [
					{
					   "name": "kube-bench",
					   "image": "aquasec/kube-bench:latest",
					   "command": [
						  "/bin/sh",
						  "-c"
					   ],
					   "args": [
						  "rm -rf /tmp/kuberesults/kube/$MY_HOSTNAME.out && cd /opt/kube-bench && kube-bench node --version ocp-3.11 |  egrep 'FAIL|PASS' | head -n -2   >> /tmp/kuberesults/kube/$MY_HOSTNAME.out && sed -i 's/\\]/\\]|/g' /tmp/kuberesults/kube/$MY_HOSTNAME.out && sed -i 's/$/|/g' /tmp/kuberesults/kube/$MY_HOSTNAME.out"
					   ],
					   "env": [
						  {
							 "name": "MY_HOSTNAME",
							 "valueFrom": {
								"fieldRef": {
								   "fieldPath": "spec.nodeName"
								}
							 }
						  }
					   ],
					   "securityContext": {
						  "privileged": true
					   },
					   "volumeMounts": [
						  {
							 "name": "var-lib-etcd",
							 "mountPath": "/var/lib/etcd",
							 "readOnly": true
						  },
						  {
							 "name": "var-lib-kubelet",
							 "mountPath": "/var/lib/kubelet",
							 "readOnly": true
						  },
						  {
							 "name": "etc-systemd",
							 "mountPath": "/etc/systemd",
							 "readOnly": true
						  },
						  {
							 "name": "etc-kubernetes",
							 "mountPath": "/etc/kubernetes",
							 "readOnly": true
						  },
						  {
							 "name": "usr-bin",
							 "mountPath": "/usr/local/mount-from-host/bin",
							 "readOnly": true
						  },
						  {
							 "name": "etc-origin",
							 "mountPath": "/etc/origin",
							 "readOnly": true
						  },
						  {
							 "name": "kube-storage",
							 "mountPath": "/tmp/kuberesults",
							 "readOnly": false
						  }
					   ]
					}
				 ],
				 "restartPolicy": "Always",
				 "volumes": [
					{
					   "name": "var-lib-etcd",
					   "hostPath": {
						  "path": "/var/lib/etcd"
					   }
					},
					{
					   "name": "var-lib-kubelet",
					   "hostPath": {
						  "path": "/var/lib/kubelet"
					   }
					},
					{
					   "name": "etc-systemd",
					   "hostPath": {
						  "path": "/etc/systemd"
					   }
					},
					{
					   "name": "etc-kubernetes",
					   "hostPath": {
						  "path": "/etc/kubernetes"
					   }
					},
					{
					   "name": "usr-bin",
					   "hostPath": {
						  "path": "/usr/bin"
					   }
					},
					{
					   "name": "etc-origin",
					   "hostPath": {
						  "path": "/etc/origin"
					   }
					},
					{
					   "name": "kube-storage",
					   "persistentVolumeClaim": {
						  "claimName": "kube-pvc"
					   }
					}
				 ]
			  }
		   }
		}
	 }`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authentication", "Bearer: rnN1L-E3PZHtJ9BRFA5N_7DRX0oRYsWrUgO98")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
func kubemaster() {
	url := "https://console-int.okd.local:443/apis/extensions/v1beta1/namespaces/kube-bench/daemonsets"
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLWJlbmNoIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImdvbGFuZy1zYS10b2tlbi1kYjZ4bSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJnb2xhbmctc2EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIxZjZlMjE1YS01NDE5LTExZWItOGM0MC0wMDBjMjk0M2FjNzYiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1iZW5jaDpnb2xhbmctc2EifQ.o4h_cEQzLksXCq58J2EdDogj95qjfxQjI9mNDyXFP1GJhTwbyM8CQQmiT77ic0vb4SRGjHSh_rZyn_FkbFi4NiUAd_K4aXmlv6EOyh5moBlauw_2g7ezBx_MJ-QZVdFfbvsyqICG70KlCBnbe6ARu3bHquWRPSQKx0wvCJpjwinUJHPZrlbhO7VIezAG4xHjjUdwH2v-sduw0z6OARX4oeV4XGZso6L2-O1f_1Dx67Lx22X9cLNe3_Iox14w44EDTpkbiViruZHrgwJ3zoehPIDNO_ZOR-40aCP1ineSqgLkxQQ9re4_kLo4AmT8Gv1M5CahOhDPBPEK2plBFR-C4A"
	var bearer = "Bearer " + token
	var jsonStr = []byte(`{
		"apiVersion": "extensions/v1beta1",
		"kind": "DaemonSet",
		"metadata": {
		   "name": "kube-bench-master"
		},
		"spec": {
		   "selector": {
			  "matchLabels": {
				 "app": "kube-bench-master"
			  }
		   },
		   "template": {
			  "metadata": {
				 "labels": {
					"app": "kube-bench-master"
				 }
			  },
			  "spec": {
				 "hostPID": true,
				 "nodeSelector": {
					"node-role.kubernetes.io/master": "true"
				 },
				 "containers": [
					{
					   "name": "kube-bench",
					   "image": "aquasec/kube-bench:latest",
					   "command": [
						  "/bin/sh",
						  "-c"
					   ],
					   "args": [
						  "rm -rf /tmp/kuberesults/kube/$MY_HOSTNAME-master.out && cd /opt/kube-bench && kube-bench master --version ocp-3.11 |   egrep 'FAIL|PASS' | head -n -2   >> /tmp/kuberesults/kube/$MY_HOSTNAME-master.out && sed -i 's/\\]/\\]|/g' /tmp/kuberesults/kube/$MY_HOSTNAME-master.out && sed -i 's/$/|/g' /tmp/kuberesults/kube/$MY_HOSTNAME-master.out"
					   ],
					   "env": [
						  {
							 "name": "MY_HOSTNAME",
							 "valueFrom": {
								"fieldRef": {
								   "fieldPath": "spec.nodeName"
								}
							 }
						  }
					   ],
					   "securityContext": {
						  "privileged": true
					   },
					   "volumeMounts": [
						  {
							 "name": "var-lib-etcd",
							 "mountPath": "/var/lib/etcd",
							 "readOnly": true
						  },
						  {
							 "name": "var-lib-kubelet",
							 "mountPath": "/var/lib/kubelet",
							 "readOnly": true
						  },
						  {
							 "name": "etc-systemd",
							 "mountPath": "/etc/systemd",
							 "readOnly": true
						  },
						  {
							 "name": "etc-kubernetes",
							 "mountPath": "/etc/kubernetes",
							 "readOnly": true
						  },
						  {
							 "name": "usr-bin",
							 "mountPath": "/usr/local/mount-from-host/bin",
							 "readOnly": true
						  },
						  {
							 "name": "etc-origin",
							 "mountPath": "/etc/origin",
							 "readOnly": true
						  },
						  {
							 "name": "kube-storage",
							 "mountPath": "/tmp/kuberesults",
							 "readOnly": false
						  }
					   ]
					}
				 ],
				 "restartPolicy": "Always",
				 "volumes": [
					{
					   "name": "var-lib-etcd",
					   "hostPath": {
						  "path": "/var/lib/etcd"
					   }
					},
					{
					   "name": "var-lib-kubelet",
					   "hostPath": {
						  "path": "/var/lib/kubelet"
					   }
					},
					{
					   "name": "etc-systemd",
					   "hostPath": {
						  "path": "/etc/systemd"
					   }
					},
					{
					   "name": "etc-kubernetes",
					   "hostPath": {
						  "path": "/etc/kubernetes"
					   }
					},
					{
					   "name": "usr-bin",
					   "hostPath": {
						  "path": "/usr/bin"
					   }
					},
					{
					   "name": "etc-origin",
					   "hostPath": {
						  "path": "/etc/origin"
					   }
					},
					{
					   "name": "kube-storage",
					   "persistentVolumeClaim": {
						  "claimName": "kube-pvc"
					   }
					}
				 ]
			  }
		   }
		}
	 }`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authentication", "Bearer: rnN1L-E3PZHtJ9BRFA5N_7DRX0oRYsWrUgO98")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
func kubeDel() {
	url := "https://console-int.okd.local:443/apis/extensions/v1beta1/namespaces/kube-bench/daemonsets/kube-bench"
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLWJlbmNoIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImdvbGFuZy1zYS10b2tlbi1kYjZ4bSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJnb2xhbmctc2EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIxZjZlMjE1YS01NDE5LTExZWItOGM0MC0wMDBjMjk0M2FjNzYiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1iZW5jaDpnb2xhbmctc2EifQ.o4h_cEQzLksXCq58J2EdDogj95qjfxQjI9mNDyXFP1GJhTwbyM8CQQmiT77ic0vb4SRGjHSh_rZyn_FkbFi4NiUAd_K4aXmlv6EOyh5moBlauw_2g7ezBx_MJ-QZVdFfbvsyqICG70KlCBnbe6ARu3bHquWRPSQKx0wvCJpjwinUJHPZrlbhO7VIezAG4xHjjUdwH2v-sduw0z6OARX4oeV4XGZso6L2-O1f_1Dx67Lx22X9cLNe3_Iox14w44EDTpkbiViruZHrgwJ3zoehPIDNO_ZOR-40aCP1ineSqgLkxQQ9re4_kLo4AmT8Gv1M5CahOhDPBPEK2plBFR-C4A"
	var bearer = "Bearer " + token
	var jsonStr = []byte(`{"kind":"DeleteOptions","apiVersion":"v1","propagationPolicy":"Foreground"}`)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authentication", "Bearer: rnN1L-E3PZHtJ9BRFA5N_7DRX0oRYsWrUgO98")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
func kubeDel1() {
	url := "https://console-int.okd.local:443/apis/extensions/v1beta1/namespaces/kube-bench/daemonsets/kube-bench-master"
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLWJlbmNoIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImdvbGFuZy1zYS10b2tlbi1kYjZ4bSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJnb2xhbmctc2EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIxZjZlMjE1YS01NDE5LTExZWItOGM0MC0wMDBjMjk0M2FjNzYiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1iZW5jaDpnb2xhbmctc2EifQ.o4h_cEQzLksXCq58J2EdDogj95qjfxQjI9mNDyXFP1GJhTwbyM8CQQmiT77ic0vb4SRGjHSh_rZyn_FkbFi4NiUAd_K4aXmlv6EOyh5moBlauw_2g7ezBx_MJ-QZVdFfbvsyqICG70KlCBnbe6ARu3bHquWRPSQKx0wvCJpjwinUJHPZrlbhO7VIezAG4xHjjUdwH2v-sduw0z6OARX4oeV4XGZso6L2-O1f_1Dx67Lx22X9cLNe3_Iox14w44EDTpkbiViruZHrgwJ3zoehPIDNO_ZOR-40aCP1ineSqgLkxQQ9re4_kLo4AmT8Gv1M5CahOhDPBPEK2plBFR-C4A"
	var bearer = "Bearer " + token
	var jsonStr = []byte(`{"kind":"DeleteOptions","apiVersion":"v1","propagationPolicy":"Foreground"}`)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authentication", "Bearer: rnN1L-E3PZHtJ9BRFA5N_7DRX0oRYsWrUgO98")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func getResults(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	fileName := string(reqBody)
	if err != nil {
		log.Fatal(err)
	}
	var displayArr []Result
	data, err := ioutil.ReadFile(root + fileName)
	if err != nil {
		fmt.Println("File reading error", err)
	}
	responses := strings.Split(string(data), "[")
	for _, resp := range responses[1:] {
		res := strings.Split(resp, "]")
		display := Result{
			Result: res[0],
			Desc:   strings.TrimSpace(res[1]),
		}
		displayArr = append(displayArr, display)
	}
	fmt.Println(displayArr)
	json, _ := json.Marshal(displayArr)
	fmt.Println(string(json))
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}
func homePage1(w http.ResponseWriter, r *http.Request) {
	kubeCall()
	kubemaster()
	time.Sleep(20 * time.Second)
	kubeDel()
	kubeDel1()
	fileArray := getFiles()
	json, _ := json.Marshal(fileArray)
	fmt.Println(string(json))
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
	// w.Write([]byte(string(json)))
}
func homePage(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}
func handleRequests() {
	http.HandleFunc("/results", getResults)
	http.HandleFunc("/scan", homePage1)
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
func getFiles() []File {
	// root := "./data"
	var filesArr []File
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		scan := File{
			Name: file.Name(),
		}
		filesArr = append(filesArr, scan)
		fmt.Println(filesArr)
	}
	return filesArr
}

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")))) //Go looks in the relative "static" directory first using http.FileServer(), then matches it to a
	//url of our choice as shown in http.Handle("/static/"). This url is what we need when referencing our css files
	//once the server begins. Our html code would therefore be <link rel="stylesheet"  href="/static/stylesheet/...">
	//It is important to note the url in http.Handle can be whatever we like, so long as we are consistent.
	handleRequests()
}
