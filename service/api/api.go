package api

// Several import statements. Importing various packages, including HTTP-related packages as well as other libraries and services.
import (
	"context"
	"time"

	goHttp "net/http"

	"github.com/pterm/pterm"
	httpIface "github.com/taubyte/http"
	http "github.com/taubyte/http/basic"
	"github.com/taubyte/http/options"
	"github.com/taubyte/tau/libdream/common"
	"github.com/taubyte/tau/libdream/services"
)

// This code defines a struct called multiverseService. It embeds a field rest of type httpIface.Service,
//  Which suggests that this struct is intended to be used as a service for handling HTTP requests.
type multiverseService struct {
	rest httpIface.Service
	common.Multiverse
}

// The BigBang function is the entry point of this code. It returns an error if something goes wrong during its execution.
func BigBang() error {
	var err error

// It initializes a multiverseService instance named srv.
// Within the initialization, it creates a Multiverse object using services.NewMultiVerse().
	srv := &multiverseService{
		Multiverse: services.NewMultiVerse(),
	}

// This section configures and creates an HTTP server using the packages.
// It sets the server to listen on a specified address (common.DreamlandApiListen) and allows all origins for (CORS) policy.
// Any errors encountered during server setup are stored in the err variable.
	srv.rest, err = http.New(srv.Context(), options.Listen(common.DreamlandApiListen), options.AllowedOrigins(true, []string{".*"}))
	if err != nil {
		return err
	}

// This code configures HTTP routes and starts the HTTP server by calling the setUpHttpRoutes method of the srv instance.
	srv.setUpHttpRoutes().Start()

// It sets a timeout for the waiting period of 10 seconds using context.WithTimeout to avoid indefinite waiting.
	waitCtx, waitCtxC := context.WithTimeout(srv.Context(), 10*time.Second)
	defer waitCtxC()

// The code enters a loop to wait for the server to become ready.
// Inside the loop, it checks for the timeout or errors and returns an error if encountered.
// It repeatedly sends HTTP GET requests to the server's address to determine if the server is ready.
// When the server becomes ready (no error in the GET request), it prints a message and returns nil, indicating a successful start.
	for {
		select {
		case <-waitCtx.Done():
			return waitCtx.Err()
		case <-time.After(100 * time.Millisecond):
			if srv.rest.Error() != nil {
				pterm.Error.Println("Dreamland failed to start")
				return srv.rest.Error()
			}
			_, err := goHttp.Get("http://" + common.DreamlandApiListen)
			if err == nil {
				pterm.Info.Println("Dreamland ready")
				return nil
			}
		}
	}
}
