package prometheus

import (
	"fmt"
	"net/http"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server ...
type Server struct{}

var (
	// TwitchMessagesIn stores the total number of messages the application has received via Twitch IRC.
	TwitchMessagesIn = prom.NewCounter(prom.CounterOpts{
		Name: "twitch_messages_received",
		Help: "Number of Twitch messages received since startup",
	})
	// TwitchMessagesOut stores the total number of messages the application has sent via Twitch IRC.
	TwitchMessagesOut = prom.NewCounter(prom.CounterOpts{
		Name: "twitch_messages_sent",
		Help: "Number of Twitch messages sent since startup",
	})

	port = 9091
)

// Init starts the prometheus server and adds the graphs.
func Init() {
	prom.MustRegister(TwitchMessagesIn)
	prom.MustRegister(TwitchMessagesOut)

	prom := &Server{}
	fmt.Printf("Prometheus server is starting on 0.0.0.0:%d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), prom); err != nil {
		fmt.Println(err)
		return
	}
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
