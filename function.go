package function

import (
	"context"
	"fmt"
	"net/http"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
)

// createHTTPTask creates a new task in your with a HTTP target.
func createHTTPTask(projectID, locationID, queueID, url, message string) (*taskspb.Task, error) {

	// Create a new Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2beta3
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %v", err)
	}

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2beta3#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2beta3#HttpRequest
			PayloadType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
				},
			},
		},
	}

	// Add a payload message if one is present.
	req.Task.GetHttpRequest().Body = []byte(message)

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}

	return createdTask, nil
}

func Execute(w http.ResponseWriter, r *http.Request) {
	_, err := createHTTPTask(os.Getenv("GOOGLE_CLOUD_PROJECT_ID"), os.Getenv("LOCATION_ID"), os.Getenv("QUEUE_ID"), os.Getenv("TARGET_URL"), "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, "request is received.")
}
