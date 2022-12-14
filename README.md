# kdeploy

Deploy from the terminal on Kubernetes.

Searches for images of requested microservice in Google Container Registry,
prompts you to interactively select an image for deployment (arrows navigation, search features),
and sets the selected image in the workload.  
> kdeploy [microservice]

If microservice was not specified - it obtains possible repositories from the registry and prompts you to select it first.  

kdeploy remembers every deployment you made and allows you to redeploy previous images.
>kdeploy --previous [microservice]

Running deploy-previous mode without specifying microservice results in prompting microservice first from your previous deployments.

### Configuration

kdeploy requires two configuration properties - `registry` and `repository`.  
The `registry` is where to look for your images (e.x. `us.gcr.io`), and the `repository` is the path to your images. If not set you will be prompted to enter them.  
Set them using:
> kdeploy config set [registry|repository] [value]

Or edit configuration file manually: 
> kdeploy config edit

Assumed that all workloads are of Deployment type. If some are StatefulSets, set them in configurations (comma separated):  
>kdeploy config set statefulsets ms-events,ms-core

### How to make it run?

Run `go build` command in the directory to build the binary file.
Then place it on at `/usr/local/bin` for convenient access anywhere on your system.
