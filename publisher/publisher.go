package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lixin9311/zapx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/sha3"
)

const (
	sigHeader = "x-signature"
)

var (
	apiKey    = os.Getenv("API_KEY")     // API key from moralis dashboard
	projectID = os.Getenv("GCP_PROJECT") // GCP project name
	deadTopic = os.Getenv("DEAD_TOPIC")  // publisher-dead

	client *pubsub.Client
)

type DeadMsg struct {
	RequestPath string `json:"request_path"`
	Data        []byte `json:"data"`
}

func init() {
	logger := zapx.Zap(zapcore.InfoLevel, zapx.WithProjectID(projectID), zapx.WithService(os.Getenv("K_SERVICE")), zapx.WithVersion(os.Getenv("K_VERSION")))
	zap.ReplaceGlobals(logger)
	// err is pre-declared to avoid shadowing client.
	if deadTopic == "" {
		zap.L().Fatal("must specify dead topic")
	}
	var err error

	// client is initialized with context.Background() because it should
	// persist between function invocations.
	client, err = pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		zap.L().Fatal("failed to create pubsub client", zap.Error(err))
	}
}

func Publish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// check API key
	sig := r.Header.Get(sigHeader)
	if sig == "" {
		zap.L().Warn("unauthorized request", zap.Any("header", r.Header))
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "unauthorized\n")
		return
	}

	// read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		zap.L().Error("failed to read body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to read body: %v\n", err)
		return
	}

	// verify alchemy signature
	if !IsValidSignatureForStringBody(body, sig, []byte(apiKey)) {
		zap.L().Warn("unauthorized request", zap.Any("header", r.Header))
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "unauthorized\n")
		return
	}

	var topic string
	// for a cloud function, the URL will be like https://asia-southeast1-runblox-openblox.cloudfunctions.net/runblox-transfer/extra-things
	path := r.URL.Path
	path = strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")
	segs := strings.Split(path, "/")
	if len(segs) > 0 {
		topic = segs[0]
	}

	if topic == "" {
		zap.L().Error("no topic name specified, push to dead topic", zap.String("path", r.URL.Path), zap.String("dead_topic", deadTopic))
		if id, err := pushToDead(ctx, body, r.URL.Path); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed to publish message: %v\n", err)
			return
		} else {
			fmt.Fprint(w, html.EscapeString("published message with id: "+id+"\n"))
			return
		}
	}

	id, err := client.Topic(topic).Publish(ctx, &pubsub.Message{
		Data: body,
	}).Get(ctx)
	if err != nil {
		if deadTopic != "" {
			zap.L().Error("failed to publish message, try to push to dead topic", zap.Error(err), zap.String("topic", topic), zap.String("dead_topic", deadTopic))
			if id, err := pushToDead(ctx, body, r.URL.Path); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "failed to publish message: %v\n", err)
				return
			} else {
				fmt.Fprint(w, html.EscapeString("published message with id: "+id+"\n"))
				return
			}
		}
		zap.L().Error("failed to publish message", zap.Error(err), zap.String("topic", topic))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to publish message: %v\n", err)
		return
	}
	fmt.Fprint(w, html.EscapeString("published message with id: "+id+"\n"))
}

func pushToDead(ctx context.Context, data []byte, path string) (string, error) {
	msg, _ := json.Marshal(&DeadMsg{RequestPath: path, Data: data})
	id, err := client.Topic(deadTopic).Publish(ctx, &pubsub.Message{
		Data: msg,
	}).Get(ctx)
	if err != nil {
		zap.L().Error("failed to publish message to dead topic",
			zap.Error(err),
			zap.String("dead_topic", deadTopic),
			zap.ByteString("data", data), zap.String("path", path),
		)
		return "", err
	}
	return id, nil
}

func IsValidSignatureForStringBody(
	body []byte, // must be raw string body, not json transformed version of the body
	signature string, // your "X-Alchemy-Signature" from header
	signingKey []byte, // taken from dashboard for specific webhook
) bool {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(body)
	hasher.Write([]byte(apiKey))
	hash := hasher.Sum(nil)
	return hexutil.Encode(hash) == signature
}
