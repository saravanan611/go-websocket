### Documentation: WebSocket Server with Gin Framework

#### **Overview**
This application demonstrates a simple WebSocket server implemented using the **Gin framework** and **Gorilla WebSocket** library. The server supports handling WebSocket connections and sending data to connected clients.

---

#### **Key Components**

1. **Gin Router**:
   - Acts as the HTTP server and router for handling routes.
   - Serves as an entry point for WebSocket requests.

2. **WebSocket Server**:
   - Manages WebSocket connections and handles data publishing.
   - Provides endpoints to connect clients and send messages.

3. **Routes**:
   - `/GetConnection`: Endpoint for clients to establish WebSocket connections.
   - `/SendData`: Endpoint to send data to connected clients.

---

#### **Code Explanation**

##### **Imports**
```go
"fmt"
"log"
_WSUtil "websocket/apps/websocket" // WebSocket utility package
"websocket/base"                  // Base configuration and constants
"github.com/gin-gonic/gin"        // Gin framework
"net/http"
```

- `_WSUtil`: Contains the WebSocket server logic.
- `base`: Includes constants like `GBaseAPI` (API base path) and `GPortNo` (server port).

---

##### **Main Function**
The `main()` function initializes the server and defines the routes.

1. **Router Initialization**
   ```go
   r := gin.Default()
   ```
   - Initializes a Gin router with default middleware.

2. **WebSocket Server Initialization**
   ```go
   wsserver := _WSUtil.NewWebSocketServer("*")
   ```
   - Creates a new WebSocket server instance allowing connections from all origins (`*`).

3. **Routes Definition**
   - **Establish WebSocket Connection**:
     ```go
     r.GET(base.GBaseAPI+"/GetConnection", func(c *gin.Context) {
         wsserver.GetConnection(c.Writer, c.Request)
     })
     ```
     - Handles WebSocket connection requests.

   - **Send Data**:
     ```go
     r.POST(base.GBaseAPI+"/SendData", func(c *gin.Context) {
         wsserver.GetData(c.Writer, c.Request)
     })
     ```
     - Publishes data to connected clients.

   - **Fallback Route**:
     ```go
     r.NoRoute(func(c *gin.Context) {
         c.Redirect(http.StatusMovedPermanently, base.GBaseAPI+"/")
     })
     ```
     - Redirects invalid or undefined routes to the base API path.

4. **Start the Server**
   ```go
   log.Printf("WebSocket server is running at :%s...", base.GPortNo)
   if err := r.Run(fmt.Sprintf(":%s", base.GPortNo)); err != nil {
       log.Fatalf("server error: %v", err)
   }
   ```
   - Starts the Gin HTTP server on the specified port (`base.GPortNo`).

---

#### **Configuration**

1. **Base Configuration** (`websocket/base`):
   - `GBaseAPI`: Base API path for routing.
   - `GPortNo`: Server port number.

---

#### **Endpoints**

| Method | Endpoint                  | Description                               |
|--------|---------------------------|-------------------------------------------|
| `GET`  | `/GetConnection`          | Establishes a WebSocket connection.       |
| `POST` | `/SendData`               | Sends data to connected WebSocket clients.|
| `ANY`  | (Undefined routes)        | Redirects to the base API path.           |

---

#### **Usage**

1. **Start the Server**:
   - Run the application using `go run main.go`.
   - The server starts listening on the port defined in `base.GPortNo`.

2. **Connect Clients**:
   - Use a WebSocket client (e.g., browser, Postman) to connect to `/GetConnection`.

3. **Send Data**:
   - Use an HTTP client (e.g., Postman) to send `POST` requests to `/SendData` with the required payload.

---

#### **Logs**

- Server logs display connection statuses, errors, and incoming requests:
  ```
  WebSocket server is running at :8080...
  ```

---

#### **Future Improvements**

1. Add authentication for WebSocket connections.
2. Implement connection tracking and metrics.
3. Enhance error handling and logging.

This documentation serves as a quick guide for understanding and using the WebSocket server.
