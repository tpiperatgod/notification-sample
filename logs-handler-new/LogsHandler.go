package logshandler

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	ofctx "github.com/OpenFunction/functions-framework-go/openfunction-context"
	alert "github.com/prometheus/alertmanager/template"
)

const (
	HTTPCodeNotFound = "404"
	Namespace        = "demo-project"
	PodName          = "wordpress-v1-f54f697c5-hdn2z"
	AlertName        = "404 Request"
	Severity         = "warning"
)

func LogsHandler(ctx *ofctx.OpenFunctionContext, in []byte) int {
	content := string(in)
	matchHTTPCode, _ := regexp.MatchString(fmt.Sprintf(" %s ", HTTPCodeNotFound), content)
	matchNamespace, _ := regexp.MatchString(fmt.Sprintf("namespace_name\":\"%s", Namespace), content)
	matchPodName, _ := regexp.MatchString(fmt.Sprintf("pod_name\":\"%s", PodName), content)

	if matchHTTPCode && matchNamespace && matchPodName {
		log.Printf("Match log - Content: %s", content)

		re := regexp.MustCompile(`([A-Z]+) (/\S*) HTTP`)
		match := re.FindAllStringSubmatch(content, -1)[0]
		path := match[len(match)-1]
		method := match[len(match)-2]

		notify := &alert.Data{
			Receiver:          "notification_manager",
			Status:            "firing",
			Alerts:            alert.Alerts{},
			GroupLabels:       alert.KV{"alertname": AlertName, "namespace": Namespace},
			CommonLabels:      alert.KV{"alertname": AlertName, "namespace": Namespace, "severity": Severity},
			CommonAnnotations: alert.KV{},
			ExternalURL:       "",
		}
		alt := alert.Alert{
			Status: "firing",
			Labels: alert.KV{
				"alertname": AlertName,
				"namespace": Namespace,
				"severity":  Severity,
				"pod":       PodName,
				"path":      path,
				"method":    method,
			},
			Annotations:  alert.KV{},
			StartsAt:     time.Now(),
			EndsAt:       time.Time{},
			GeneratorURL: "",
			Fingerprint:  "",
		}
		notify.Alerts = append(notify.Alerts, alt)
		notifyBytes, _ := json.Marshal(notify)
		if err := ctx.SendTo(notifyBytes, "alert"); err != nil {
			panic(err)
		}
		log.Printf("Send log to notification manager.")
	}
	return 200
}
