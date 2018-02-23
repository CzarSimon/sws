package proxy

import (
	"io/ioutil"
	"testing"

	"github.com/CzarSimon/sws/pkg/service"
)

const testPath = "test/nginx.conf"
const expectedConf = "worker_processes 1;\n\npid /var/run/nginx.pid;\n\nevents {\n\tworker_connections 1024;\n}\n\nhttp {\n\tinclude mime.types;\n\tdefault_type application/octet-stream;\n\tsendfile on;\n\n\tserver {\n\t\tlisten 18080;\n\n\t\tlocation /example.com/ {\n\t\t\tproxy_pass http://example-service:8080/;\n\t\t}\n\n\t\tlocation /second-example.com/ {\n\t\t\tproxy_pass http://second-service:7070/;\n\t\t}\n\n\t\tlocation = /sws-proxy/health {\n\t\t\treturn 200;\n\t\t\taccess_log off;\n\t\t}\n\t}\n}"

func TestCreateConf(t *testing.T) {
	p := New("sws-proxy-1", 18080, getTestServices())
	conf := p.CreateConf("/sws-proxy/health")
	if conf != expectedConf {
		t.Errorf("CreateConf Error!\nExpected:\n%s\nGot:\n%s", expectedConf, conf)
	}
}

func TestWriteConf(t *testing.T) {
	p := New("sws-proxy-1", 18080, getTestServices())
	conf := p.CreateConf("/sws-proxy/health")
	err := WriteConf(conf, testPath)
	if err != nil {
		t.Fatalf("WriteConf returned an error:\n%s", err.Error())
	}
	bytes, err := ioutil.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Reading config returned an error:\n%s", err.Error())
	}
	recievedConf := string(bytes)
	if recievedConf != expectedConf {
		t.Errorf("WriteConf Error!\nExpected:\n%s\nGot:\n%s", expectedConf, recievedConf)
	}
}

func getTestServices() []service.Service {
	s1 := getTestService()
	s2 := getTestService()
	s2.Domain = "second-example.com"
	s2.Name = "second-service"
	s2.Port = 7070
	return []service.Service{s1, s2}
}
