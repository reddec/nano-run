package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"nano-run/services/meta"
	"nano-run/worker"
)

// Expose worker as HTTP handler:
//     POST /                  - post task async, returns 303 See Other and location.
//     GET  /:id               - get task info.
//     GET  /:id/completed     - redirect to completed attempt (or 404 if task not yet complete)
//     GET  /:id/attempt/:atid - get attempt result (as-is).
//     GET  /:id/request       - replay request (as-is).
func Expose(router *mux.Router, wrk *worker.Worker) {
	//TODO: wait
	router.Path("/").Methods("POST").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id, err := wrk.Enqueue(request)
		if err != nil {
			log.Println("failed to enqueue:", err)
			http.Error(writer, "failed to enqueue", http.StatusInternalServerError)
			return
		}
		writer.Header().Set("X-Correlation-Id", id)
		http.Redirect(writer, request, id, http.StatusSeeOther)
	})
	// get state: 200 with json description. For complete request - Location header will be filled.
	router.Path("/{id}").Methods("GET").HandlerFunc(createTask(wrk))
	router.Path("/{id}/completed").Methods("GET").HandlerFunc(getComplete(wrk))
	// get attempt result as-is.
	router.Path("/{id}/attempt/{attemptId}").Methods("GET").HandlerFunc(getAttempt(wrk))
	// get recorded request.
	router.Path("/{id}/request").Methods("GET").HandlerFunc(getRequest(wrk))
}

// get request meta information.
func createTask(wrk *worker.Worker) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)
		requestID := params["id"]
		info, err := wrk.Meta().Get(requestID)
		if err != nil {
			log.Println("failed access request", requestID, ":", err)
			http.NotFound(writer, request)
			return
		}
		data, err := json.Marshal(info)
		if err != nil {
			log.Println("failed to encode info:", err)
			http.Error(writer, "encoding", http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Content-Length", strconv.Itoa(len(data)))
		writer.Header().Set("Content-Version", strconv.Itoa(len(info.Attempts)))
		// modification time
		setLastModify(writer, info)

		writer.Header().Set("Age", strconv.FormatInt(int64(time.Since(info.CreatedAt)/time.Second), 10))
		if info.Complete {
			writer.Header().Set("X-Status", "complete")
		} else {
			writer.Header().Set("X-Status", "processing")
		}

		if len(info.Attempts) > 0 {
			writer.Header().Set("X-Last-Attempt", info.Attempts[len(info.Attempts)-1].ID)
			writer.Header().Set("X-Last-Attempt-At", info.Attempts[len(info.Attempts)-1].CreatedAt.Format(time.RFC850))
		}

		if info.Complete {
			lastAttempt := info.Attempts[len(info.Attempts)-1]
			request.URL.Path += "/attempt/" + lastAttempt.ID
			writer.Header().Set("Location", request.URL.String())
		}
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write(data)
	}
}

func getComplete(wrk *worker.Worker) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)
		requestID := params["id"]
		info, err := wrk.Meta().Get(requestID)
		if err != nil {
			log.Println("failed access request", requestID, ":", err)
			http.NotFound(writer, request)
			return
		}
		if !info.Complete {
			http.NotFound(writer, request)
			return
		}
		lastAttempt := info.Attempts[len(info.Attempts)-1]
		http.Redirect(writer, request, "attempt/"+lastAttempt.ID, http.StatusMovedPermanently)
	}
}

func getAttempt(wrk *worker.Worker) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)
		requestID := params["id"]
		attemptID := params["attemptId"]
		info, err := wrk.Meta().Get(requestID)
		if err != nil {
			log.Println("failed access request", requestID, ":", err)
			http.NotFound(writer, request)
			return
		}
		var attempt meta.Attempt
		var found bool
		for _, atp := range info.Attempts {
			if atp.ID == attemptID {
				found = true
				attempt = atp
				break
			}
		}
		if !found {
			http.NotFound(writer, request)
			return
		}
		body, err := wrk.Blobs().Get(attempt.ID)
		if err != nil {
			log.Println("failed to get body:", err)
			http.Error(writer, "get body", http.StatusInternalServerError)
			return
		}
		defer body.Close()
		writer.Header().Set("Last-Modified", attempt.CreatedAt.Format(time.RFC850))
		if info.Complete {
			writer.Header().Set("X-Status", "complete")
		} else {
			writer.Header().Set("X-Status", "processing")
		}
		writer.Header().Set("X-Processed", "true")
		for k, v := range attempt.Headers {
			writer.Header()[k] = v
		}
		writer.WriteHeader(attempt.Code)
		_, _ = io.Copy(writer, body)
	}
}

func getRequest(wrk *worker.Worker) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)
		requestID := params["id"]
		info, err := wrk.Meta().Get(requestID)
		if err != nil {
			log.Println("failed access request", requestID, ":", err)
			http.NotFound(writer, request)
			return
		}
		writer.Header().Set("Last-Modified", info.CreatedAt.Format(time.RFC850))
		f, err := wrk.Blobs().Get(requestID)
		if err != nil {
			log.Println("failed to get data:", err)
			http.Error(writer, "data", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		writer.Header().Set("X-Method", info.Method)
		writer.Header().Set("X-Request-Uri", info.URI)

		for k, v := range info.Headers {
			writer.Header()[k] = v
		}
		writer.WriteHeader(http.StatusOK)
		_, _ = io.Copy(writer, f)
	}
}

func setLastModify(writer http.ResponseWriter, info *meta.Request) {
	if info.Complete {
		writer.Header().Set("Last-Modified", info.CompleteAt.Format(time.RFC850))
	} else if len(info.Attempts) > 0 {
		writer.Header().Set("Last-Modified", info.Attempts[len(info.Attempts)-1].CreatedAt.Format(time.RFC850))
	} else {
		writer.Header().Set("Last-Modified", info.CreatedAt.Format(time.RFC850))
	}
}
