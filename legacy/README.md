# kdeploy
Deploy from the terminal.

Searches for images of requested microservice in Google Container Registry,
prompts interactively you to select an image for deployment (arrows navigation, search features),
sets the selected image in workload, and conveniently shows you watch over deploying microservices.

Also has a deploy-previous mode - when you need quickly redeploy what was before your last deployment.
However, it has goldfish memory - can redeploy only the previous deployment.

### Synopsis
```
# kdeploy <microservice> [images-list-size]
# kdeploy previous
```

### Examples:
```
# kdeploy data-generator
# kdeploy risk-manager 50

# kdeploy previous
```

### How to make it run?

First, you need to run `go build` command in the directory to build the `go` binary (there is already a binary for OSX, you can use it).
Then you need to make `kdeploy` and `select_image_prompt` files executable and place them on at `/usr/local/bin`.

### How to use it?

Run command specifying microservice name and optionally number of images to be listed.
Tool will interactively prompt you for selecting image to deploy.
After selection, it will gather all additionally needed info as full digest of image and using kubectl update image
in deployment/statefulset.

### How it works?

It gets images list and related info with help of gcloud command, updates image with kubectl.
For interactive selection go script is used.
